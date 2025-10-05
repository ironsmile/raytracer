package film

import (
    "cmp"
    "fmt"
    "math"
    "slices"
    "time"
    "unsafe"

    "github.com/ironsmile/raytracer/camera"
    "github.com/ironsmile/raytracer/engine"
    "github.com/ironsmile/raytracer/film/shaders"
    "github.com/ironsmile/raytracer/optional"
    "github.com/ironsmile/raytracer/sampler"
    "github.com/ironsmile/raytracer/scene"
    "github.com/ironsmile/raytracer/unsafer"

    "github.com/go-gl/glfw/v3.3/glfw"
    vk "github.com/vulkan-go/vulkan"
)

type VulkanAppArgs struct {
    Fullscreen  bool
    VSync       bool
    Width       int
    Height      int
    Interactive bool
    ShowBBoxes  bool
    FPSCap      uint
    ShowFPS     bool
    SceneName   string

    // Debug causes few additional diagnostics messages to be printed while working.
    Debug bool
}

const (
    title             = "Raytracer"
    maxFramesInFlight = 2
)

func NewVulkanWindow(args VulkanAppArgs) *VulkanApp {
    return &VulkanApp{
        enableValidationLayers: args.Debug,
        validationLayers: []string{
            "VK_LAYER_KHRONOS_validation\x00",
        },
        deviceExtensions: []string{
            vk.KhrSwapchainExtensionName + "\x00",
        },
        physicalDevice: vk.PhysicalDevice(vk.NullHandle),
        device:         vk.Device(vk.NullHandle),
        surface:        vk.NullSurface,
        swapChain:      vk.NullSwapchain,
        args:           args,
    }
}

// VulkanApp is a `film` which renders the scene with Vulkan and GLFW3.
type VulkanApp struct {
    // validationLayers is the list of required device extensions needed by this
    // program when the -D flag is set.
    validationLayers       []string
    enableValidationLayers bool

    // deviceExtensions is the list of required device extensions needed by this
    // program.
    deviceExtensions []string

    window   *glfw.Window
    instance vk.Instance

    // physicalDevice is the physical device selected for this program.
    physicalDevice vk.PhysicalDevice

    // device is the logical device created for interfacing with the physical device.
    device vk.Device

    graphicsQueue vk.Queue
    presentQueue  vk.Queue

    surface vk.Surface

    swapChain            vk.Swapchain
    swapChainImages      []vk.Image
    swapChainImageViews  []vk.ImageView
    swapChainImageFormat vk.Format
    swapChainExtend      vk.Extent2D

    swapChainFramebuffers []vk.Framebuffer

    renderPass     vk.RenderPass
    pipelineLayout vk.PipelineLayout

    graphicsPipline vk.Pipeline

    commandPool    vk.CommandPool
    commandBuffers []vk.CommandBuffer

    imageAvailabmeSems []vk.Semaphore
    renderFinishedSems []vk.Semaphore
    inFlightFences     []vk.Fence

    frameBufferResized bool

    curentFrame uint32

    // Raytracer engine stuff
    args    VulkanAppArgs
    film    *vulkanFilm
    sampler *sampler.SimpleSampler
    tracer  *engine.FPSEngine
    cam     camera.Camera

    // Copying film to GPU memory stuff
    filmStagingBuffer       vk.Buffer
    filmStagingBufferMemory vk.DeviceMemory
    filmBufferData          unsafe.Pointer

    filmImageFormat vk.Format
    filmImage       vk.Image
    filmImageMemory vk.DeviceMemory
}

// Run runs the vulkan program.
func (a *VulkanApp) Run() error {
    if err := a.initWindow(); err != nil {
        return fmt.Errorf("initWindow: %w", err)
    }
    defer a.cleanWindow()

    if err := a.initVulkan(); err != nil {
        return fmt.Errorf("initVulkan: %w", err)
    }
    defer a.cleanupVulkan()

    if err := a.initEngine(); err != nil {
        return fmt.Errorf("initEngine: %w", err)
    }
    defer a.cleanEngine()

    if err := a.mainLoop(); err != nil {
        return fmt.Errorf("mainLoop: %w", err)
    }

    return nil
}

func (a *VulkanApp) initWindow() error {
    if err := glfw.Init(); err != nil {
        return fmt.Errorf("glfw.Init: %w", err)
    }

    glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

    window, err := glfw.CreateWindow(a.args.Width, a.args.Height, title, nil, nil)
    if err != nil {
        return fmt.Errorf("creating window: %w", err)
    }

    window.SetFramebufferSizeCallback(a.frameBufferResizeCallback)

    a.window = window
    return nil
}

func (a *VulkanApp) frameBufferResizeCallback(
    w *glfw.Window,
    width int,
    height int,
) {
    a.frameBufferResized = true
}

func (a *VulkanApp) cleanWindow() {
    a.window.Destroy()
    glfw.Terminate()
}

func (a *VulkanApp) initVulkan() error {
    vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())

    if err := vk.Init(); err != nil {
        return fmt.Errorf("failed to init Vulkan Go: %w", err)
    }

    if err := a.createInstance(); err != nil {
        return fmt.Errorf("createInstance: %w", err)
    }

    if err := a.createSurface(); err != nil {
        return fmt.Errorf("createSurface: %w", err)
    }

    if err := a.pickPhysicalDevice(); err != nil {
        return fmt.Errorf("pickPhysicalDevice: %w", err)
    }

    if err := a.createLogicalDevice(); err != nil {
        return fmt.Errorf("createLogicalDevice: %w", err)
    }

    if err := a.createSwapChain(); err != nil {
        return fmt.Errorf("createSwapChain: %w", err)
    }

    if err := a.createImageViews(); err != nil {
        return fmt.Errorf("createImageViews: %w", err)
    }

    if err := a.createRenderPass(); err != nil {
        return fmt.Errorf("createRenderPass: %w", err)
    }

    if err := a.createGraphicsPipeline(); err != nil {
        return fmt.Errorf("createGraphicsPipeline: %w", err)
    }

    if err := a.createFramebuffers(); err != nil {
        return fmt.Errorf("createFramebuffers: %w", err)
    }

    if err := a.createCommandPool(); err != nil {
        return fmt.Errorf("createCommandPool: %w", err)
    }

    if err := a.createCommandBuffer(); err != nil {
        return fmt.Errorf("createCommandBuffer: %w", err)
    }

    if err := a.createFilmImage(); err != nil {
        return fmt.Errorf("createFilmImage: %w", err)
    }

    if err := a.createSyncObjects(); err != nil {
        return fmt.Errorf("createSyncObjects: %w", err)
    }

    return nil
}

func (a *VulkanApp) cleanupVulkan() {
    for i := 0; i < maxFramesInFlight; i++ {
        vk.DestroySemaphore(a.device, a.imageAvailabmeSems[i], nil)
        vk.DestroySemaphore(a.device, a.renderFinishedSems[i], nil)
        vk.DestroyFence(a.device, a.inFlightFences[i], nil)
    }

    vk.DestroyCommandPool(a.device, a.commandPool, nil)

    vk.DestroyPipeline(a.device, a.graphicsPipline, nil)
    vk.DestroyPipelineLayout(a.device, a.pipelineLayout, nil)

    a.cleanupSwapChain()

    vk.DestroyRenderPass(a.device, a.renderPass, nil)

    a.cleanFilmImage()

    if a.device != vk.Device(vk.NullHandle) {
        vk.DestroyDevice(a.device, nil)
    }
    if a.surface != vk.NullSurface {
        vk.DestroySurface(a.instance, a.surface, nil)
    }
    vk.DestroyInstance(a.instance, nil)
}

func (a *VulkanApp) cleanupSwapChain() {
    for _, frameBuffer := range a.swapChainFramebuffers {
        vk.DestroyFramebuffer(a.device, frameBuffer, nil)
    }

    for _, imageView := range a.swapChainImageViews {
        vk.DestroyImageView(a.device, imageView, nil)
    }

    if a.swapChain != vk.NullSwapchain {
        vk.DestroySwapchain(a.device, a.swapChain, nil)
    }
    a.swapChainImages = nil
    a.swapChainImageViews = nil
}

func (a *VulkanApp) createSurface() error {
    surfacePtr, err := a.window.CreateWindowSurface(a.instance, nil)
    if err != nil {
        return fmt.Errorf("cannot create surface within GLFW window: %w", err)
    }

    a.surface = vk.SurfaceFromPointer(surfacePtr)
    return nil
}

func (a *VulkanApp) pickPhysicalDevice() error {
    var deviceCount uint32
    err := vk.Error(vk.EnumeratePhysicalDevices(a.instance, &deviceCount, nil))
    if err != nil {
        return fmt.Errorf("failed to get the number of physical devices: %w", err)
    }
    if deviceCount == 0 {
        return fmt.Errorf("failed to find GPUs with Vulkan support")
    }

    pDevices := make([]vk.PhysicalDevice, deviceCount)
    err = vk.Error(vk.EnumeratePhysicalDevices(a.instance, &deviceCount, pDevices))
    if err != nil {
        return fmt.Errorf("failed to enumerate the physical devices: %w", err)
    }

    var (
        selectedDevice vk.PhysicalDevice
        score          uint32
    )

    for _, device := range pDevices {
        deviceScore := a.getDeviceScore(device)

        if deviceScore > score {
            selectedDevice = device
            score = deviceScore
        }
    }

    if selectedDevice == vk.PhysicalDevice(vk.NullHandle) {
        return fmt.Errorf("failed to find suitable physical devices")
    }

    a.physicalDevice = selectedDevice
    return nil
}

func (a *VulkanApp) createLogicalDevice() error {
    indices := a.findQueueFamilies(a.physicalDevice)
    if !indices.IsComplete() {
        return fmt.Errorf("createLogicalDevice called for physical device which does " +
            "have all the queues required by the program")
    }

    queueFamilies := make(map[uint32]struct{})
    queueFamilies[indices.Graphics.Get()] = struct{}{}
    queueFamilies[indices.Present.Get()] = struct{}{}

    queueCreateInfos := []vk.DeviceQueueCreateInfo{}

    for familyIndex := range queueFamilies {
        queueCreateInfos = append(
            queueCreateInfos,
            vk.DeviceQueueCreateInfo{
                SType:            vk.StructureTypeDeviceQueueCreateInfo,
                QueueFamilyIndex: familyIndex,
                QueueCount:       1,
                PQueuePriorities: []float32{1.0},
            },
        )
    }

    //!TODO: left for later use
    deviceFeatures := []vk.PhysicalDeviceFeatures{{}}

    createInfo := vk.DeviceCreateInfo{
        SType:            vk.StructureTypeDeviceCreateInfo,
        PEnabledFeatures: deviceFeatures,

        PQueueCreateInfos:    queueCreateInfos,
        QueueCreateInfoCount: uint32(len(queueCreateInfos)),

        EnabledExtensionCount:   uint32(len(a.deviceExtensions)),
        PpEnabledExtensionNames: a.deviceExtensions,
    }

    if a.enableValidationLayers {
        createInfo.PpEnabledLayerNames = a.validationLayers
        createInfo.EnabledLayerCount = uint32(len(a.validationLayers))
    }

    var device vk.Device
    err := vk.Error(vk.CreateDevice(a.physicalDevice, &createInfo, nil, &device))
    if err != nil {
        return fmt.Errorf("failed to create logical device: %w", err)
    }
    a.device = device

    var graphicsQueue vk.Queue
    vk.GetDeviceQueue(a.device, indices.Graphics.Get(), 0, &graphicsQueue)
    a.graphicsQueue = graphicsQueue

    var presentQueue vk.Queue
    vk.GetDeviceQueue(a.device, indices.Present.Get(), 0, &presentQueue)
    a.presentQueue = presentQueue

    return nil
}

func (a *VulkanApp) createSwapChain() error {
    swapChainSupport := a.querySwapChainSupport(a.physicalDevice)

    surfaceFormat := a.chooseSwapSurfaceFormat(swapChainSupport.formats)
    presentMode := a.chooseSwapPresentMode(swapChainSupport.presentModes)
    extend := a.chooseSwapExtend(swapChainSupport.capabilities)

    imageCount := swapChainSupport.capabilities.MinImageCount + 1
    if swapChainSupport.capabilities.MaxImageCount > 0 &&
        imageCount > swapChainSupport.capabilities.MaxImageCount {
        imageCount = swapChainSupport.capabilities.MaxImageCount
    }

    createInfo := vk.SwapchainCreateInfo{
        SType:            vk.StructureTypeSwapchainCreateInfo,
        Surface:          a.surface,
        MinImageCount:    imageCount,
        ImageColorSpace:  surfaceFormat.ColorSpace,
        ImageFormat:      surfaceFormat.Format,
        ImageExtent:      extend,
        ImageArrayLayers: 1,
        ImageUsage: vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit) |
            vk.ImageUsageFlags(vk.ImageUsageTransferDstBit),
        PreTransform:   swapChainSupport.capabilities.CurrentTransform,
        CompositeAlpha: vk.CompositeAlphaOpaqueBit,
        PresentMode:    presentMode,
        Clipped:        vk.True,
    }

    if a.args.Debug {
        fmt.Printf("Selected swapchain image format: %#v, color space: %#v\n",
            createInfo.ImageFormat, createInfo.ImageColorSpace,
        )
    }

    indices := a.findQueueFamilies(a.physicalDevice)
    if indices.Graphics.Get() != indices.Present.Get() {
        createInfo.ImageSharingMode = vk.SharingModeConcurrent
        createInfo.QueueFamilyIndexCount = 2
        createInfo.PQueueFamilyIndices = []uint32{
            indices.Graphics.Get(),
            indices.Present.Get(),
        }
    } else {
        createInfo.ImageSharingMode = vk.SharingModeExclusive
    }

    var swapChain vk.Swapchain
    res := vk.CreateSwapchain(a.device, &createInfo, nil, &swapChain)
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("failed to create swap chain: %w", err)
    }
    a.swapChain = swapChain

    var imagesCount uint32
    vk.GetSwapchainImages(a.device, a.swapChain, &imagesCount, nil)

    images := make([]vk.Image, imagesCount)
    vk.GetSwapchainImages(a.device, a.swapChain, &imagesCount, images)

    a.swapChainImages = images

    a.swapChainImageFormat = surfaceFormat.Format
    a.swapChainExtend = extend

    return nil
}

func (a *VulkanApp) recreateSwapChain() error {
    for true {
        width, height := a.window.GetFramebufferSize()
        if width != 0 || height != 0 {
            break
        }

        glfw.WaitEvents()
    }

    vk.DeviceWaitIdle(a.device)

    a.cleanupSwapChain()

    if err := a.createSwapChain(); err != nil {
        return fmt.Errorf("createSwapChain: %w", err)
    }
    if err := a.createImageViews(); err != nil {
        return fmt.Errorf("createImageViews: %w", err)
    }
    if err := a.createFramebuffers(); err != nil {
        return fmt.Errorf("createFramebuffers: %w", err)
    }

    return nil
}

func (a *VulkanApp) createImageViews() error {
    for i, swapChainImage := range a.swapChainImages {
        swapChainImage := swapChainImage
        createInfo := vk.ImageViewCreateInfo{
            SType:    vk.StructureTypeImageViewCreateInfo,
            Image:    swapChainImage,
            ViewType: vk.ImageViewType2d,
            Format:   a.swapChainImageFormat,
            Components: vk.ComponentMapping{
                R: vk.ComponentSwizzleIdentity,
                G: vk.ComponentSwizzleIdentity,
                B: vk.ComponentSwizzleIdentity,
                A: vk.ComponentSwizzleIdentity,
            },
            SubresourceRange: vk.ImageSubresourceRange{
                AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
                BaseMipLevel:   0,
                LevelCount:     1,
                BaseArrayLayer: 0,
                LayerCount:     1,
            },
        }

        var imageView vk.ImageView
        res := vk.CreateImageView(a.device, &createInfo, nil, &imageView)
        if err := vk.Error(res); err != nil {
            return fmt.Errorf("failed to create image %d: %w", i, err)
        }

        a.swapChainImageViews = append(a.swapChainImageViews, imageView)
    }

    return nil
}

func (a *VulkanApp) createRenderPass() error {
    colorAttachment := vk.AttachmentDescription{
        Format:         a.swapChainImageFormat,
        Samples:        vk.SampleCount1Bit,
        LoadOp:         vk.AttachmentLoadOpClear,
        StoreOp:        vk.AttachmentStoreOpStore,
        StencilLoadOp:  vk.AttachmentLoadOpDontCare,
        StencilStoreOp: vk.AttachmentStoreOpDontCare,
        InitialLayout:  vk.ImageLayoutUndefined,
        FinalLayout:    vk.ImageLayoutPresentSrc,
    }

    colorAttachmentRef := vk.AttachmentReference{
        Attachment: 0,
        Layout:     vk.ImageLayoutColorAttachmentOptimal,
    }

    subpass := vk.SubpassDescription{
        PipelineBindPoint:    vk.PipelineBindPointGraphics,
        ColorAttachmentCount: 1,
        PColorAttachments:    []vk.AttachmentReference{colorAttachmentRef},
    }

    dependency := vk.SubpassDependency{
        SrcSubpass:    vk.SubpassExternal,
        DstSubpass:    0,
        SrcStageMask:  vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
        SrcAccessMask: 0,
        DstStageMask:  vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
        DstAccessMask: vk.AccessFlags(vk.AccessColorAttachmentWriteBit),
    }

    rederPassInfo := vk.RenderPassCreateInfo{
        SType:           vk.StructureTypeRenderPassCreateInfo,
        AttachmentCount: 1,
        PAttachments:    []vk.AttachmentDescription{colorAttachment},
        SubpassCount:    1,
        PSubpasses:      []vk.SubpassDescription{subpass},
        DependencyCount: 1,
        PDependencies:   []vk.SubpassDependency{dependency},
    }

    var renderPass vk.RenderPass
    res := vk.CreateRenderPass(a.device, &rederPassInfo, nil, &renderPass)
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("failed to create render pass: %w", err)
    }
    a.renderPass = renderPass

    return nil
}

func (a *VulkanApp) createGraphicsPipeline() error {
    vertShaderCode, err := shaders.FS.ReadFile("vert.spv")
    if err != nil {
        return fmt.Errorf("failed to read vertex shader bytecode: %w", err)
    }

    fragShaderCode, err := shaders.FS.ReadFile("frag.spv")
    if err != nil {
        return fmt.Errorf("failed to read fragment shader bytecode: %w", err)
    }

    if a.args.Debug {
        fmt.Printf("vertex shader code size: %d\n", len(vertShaderCode))
        fmt.Printf("fragment shader code size: %d\n", len(fragShaderCode))
    }

    vertexShaderModule, err := a.createShaderModule(vertShaderCode)
    if err != nil {
        return fmt.Errorf("creating vertex shader module: %w", err)
    }
    defer vk.DestroyShaderModule(a.device, vertexShaderModule, nil)

    fragmentShaderModule, err := a.createShaderModule(fragShaderCode)
    if err != nil {
        return fmt.Errorf("creating fragment shader module: %w", err)
    }
    defer vk.DestroyShaderModule(a.device, fragmentShaderModule, nil)

    vertShaderStageInfo := vk.PipelineShaderStageCreateInfo{
        SType:  vk.StructureTypePipelineShaderStageCreateInfo,
        Stage:  vk.ShaderStageVertexBit,
        Module: vertexShaderModule,
        PName:  "main\x00",
    }

    fragShaderStageInfo := vk.PipelineShaderStageCreateInfo{
        SType:  vk.StructureTypePipelineShaderStageCreateInfo,
        Stage:  vk.ShaderStageFragmentBit,
        Module: fragmentShaderModule,
        PName:  "main\x00",
    }

    shaderStages := []vk.PipelineShaderStageCreateInfo{
        vertShaderStageInfo,
        fragShaderStageInfo,
    }

    vertexInputInfo := vk.PipelineVertexInputStateCreateInfo{
        SType: vk.StructureTypePipelineVertexInputStateCreateInfo,

        VertexBindingDescriptionCount:   0,
        VertexAttributeDescriptionCount: 0,
    }

    inputAssembly := vk.PipelineInputAssemblyStateCreateInfo{
        SType:                  vk.StructureTypePipelineInputAssemblyStateCreateInfo,
        Topology:               vk.PrimitiveTopologyTriangleList,
        PrimitiveRestartEnable: vk.False,
    }

    viewport := vk.Viewport{
        X:        0,
        Y:        0,
        Width:    float32(a.swapChainExtend.Width),
        Height:   float32(a.swapChainExtend.Height),
        MinDepth: 0,
        MaxDepth: 1,
    }

    scissor := vk.Rect2D{
        Offset: vk.Offset2D{X: 0, Y: 0},
        Extent: a.swapChainExtend,
    }

    dynamicStates := []vk.DynamicState{
        vk.DynamicStateViewport,
        vk.DynamicStateScissor,
    }

    dynamicState := vk.PipelineDynamicStateCreateInfo{
        SType:             vk.StructureTypePipelineDynamicStateCreateInfo,
        DynamicStateCount: uint32(len(dynamicStates)),
        PDynamicStates:    dynamicStates,
    }

    viewportState := vk.PipelineViewportStateCreateInfo{
        SType:         vk.StructureTypePipelineViewportStateCreateInfo,
        ViewportCount: 1,
        ScissorCount:  1,
        PViewports:    []vk.Viewport{viewport},
        PScissors:     []vk.Rect2D{scissor},
    }

    rasterizer := vk.PipelineRasterizationStateCreateInfo{
        SType:                   vk.StructureTypePipelineRasterizationStateCreateInfo,
        DepthClampEnable:        vk.False,
        RasterizerDiscardEnable: vk.False,
        PolygonMode:             vk.PolygonModeFill,
        LineWidth:               1,
        CullMode:                vk.CullModeFlags(vk.CullModeBackBit),
        FrontFace:               vk.FrontFaceClockwise,
        DepthBiasEnable:         vk.False,
    }

    multisampling := vk.PipelineMultisampleStateCreateInfo{
        SType:                 vk.StructureTypePipelineMultisampleStateCreateInfo,
        SampleShadingEnable:   vk.False,
        RasterizationSamples:  vk.SampleCount1Bit,
        MinSampleShading:      1,
        AlphaToCoverageEnable: vk.False,
        AlphaToOneEnable:      vk.False,
    }

    colorBlnedAttachment := vk.PipelineColorBlendAttachmentState{
        ColorWriteMask: vk.ColorComponentFlags(
            vk.ColorComponentRBit |
                vk.ColorComponentGBit |
                vk.ColorComponentBBit |
                vk.ColorComponentABit,
        ),
        BlendEnable:         vk.False,
        SrcColorBlendFactor: vk.BlendFactorOne,
        DstColorBlendFactor: vk.BlendFactorZero,
        ColorBlendOp:        vk.BlendOpAdd,
        SrcAlphaBlendFactor: vk.BlendFactorOne,
        DstAlphaBlendFactor: vk.BlendFactorZero,
        AlphaBlendOp:        vk.BlendOpAdd,
    }

    colorBlending := vk.PipelineColorBlendStateCreateInfo{
        SType:           vk.StructureTypePipelineColorBlendStateCreateInfo,
        LogicOpEnable:   vk.False,
        LogicOp:         vk.LogicOpCopy,
        AttachmentCount: 1,
        PAttachments: []vk.PipelineColorBlendAttachmentState{
            colorBlnedAttachment,
        },
    }

    pipelineLayoutInfo := vk.PipelineLayoutCreateInfo{
        SType:          vk.StructureTypePipelineLayoutCreateInfo,
        SetLayoutCount: 0,
    }

    var pipelineLayout vk.PipelineLayout
    res := vk.CreatePipelineLayout(a.device, &pipelineLayoutInfo, nil, &pipelineLayout)
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("failed to create pipeline layout: %w", err)
    }
    a.pipelineLayout = pipelineLayout

    pipelineInfo := vk.GraphicsPipelineCreateInfo{
        SType:               vk.StructureTypeGraphicsPipelineCreateInfo,
        StageCount:          uint32(len(shaderStages)),
        PStages:             shaderStages,
        PVertexInputState:   &vertexInputInfo,
        PInputAssemblyState: &inputAssembly,
        PViewportState:      &viewportState,
        PRasterizationState: &rasterizer,
        PMultisampleState:   &multisampling,
        PDepthStencilState:  nil,
        PColorBlendState:    &colorBlending,
        PDynamicState:       &dynamicState,
        Layout:              a.pipelineLayout,
        RenderPass:          a.renderPass,
        Subpass:             0,
        BasePipelineHandle:  vk.Pipeline(vk.NullHandle),
        BasePipelineIndex:   -1,
    }

    pipelines := make([]vk.Pipeline, 1)
    res = vk.CreateGraphicsPipelines(
        a.device,
        vk.PipelineCache(vk.NullHandle),
        1,
        []vk.GraphicsPipelineCreateInfo{pipelineInfo},
        nil,
        pipelines,
    )
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("failed to create graphics pipeline: %w", err)
    }
    a.graphicsPipline = pipelines[0]

    return nil
}

func (a *VulkanApp) createFramebuffers() error {
    a.swapChainFramebuffers = make([]vk.Framebuffer, len(a.swapChainImageViews))

    for i, swapChainView := range a.swapChainImageViews {
        swapChainView := swapChainView

        attachments := []vk.ImageView{
            swapChainView,
        }

        frameBufferInfo := vk.FramebufferCreateInfo{
            SType:           vk.StructureTypeFramebufferCreateInfo,
            RenderPass:      a.renderPass,
            AttachmentCount: 1,
            PAttachments:    attachments,
            Width:           a.swapChainExtend.Width,
            Height:          a.swapChainExtend.Height,
            Layers:          1,
        }

        var frameBuffer vk.Framebuffer
        res := vk.CreateFramebuffer(a.device, &frameBufferInfo, nil, &frameBuffer)
        if err := vk.Error(res); err != nil {
            return fmt.Errorf("failed to create frame buffer %d: %w", i, err)
        }

        a.swapChainFramebuffers[i] = frameBuffer
    }

    return nil
}

func (a *VulkanApp) copyBuffer(
    srcBuffer vk.Buffer,
    dstBuffer vk.Buffer,
    size vk.DeviceSize,
) error {
    commandBuffer, err := a.beginSingleTimeCommands()
    if err != nil {
        return fmt.Errorf("failed to begin single time commands: %w", err)
    }

    copyRegion := vk.BufferCopy{
        SrcOffset: 0,
        DstOffset: 0,
        Size:      size,
    }

    vk.CmdCopyBuffer(commandBuffer, srcBuffer, dstBuffer, 1, []vk.BufferCopy{copyRegion})

    return a.endSingleTimeCommands(commandBuffer)
}

func (a *VulkanApp) beginSingleTimeCommands() (vk.CommandBuffer, error) {
    allocInfo := vk.CommandBufferAllocateInfo{
        SType:              vk.StructureTypeCommandBufferAllocateInfo,
        Level:              vk.CommandBufferLevelPrimary,
        CommandPool:        a.commandPool,
        CommandBufferCount: 1,
    }

    commandBuffers := make([]vk.CommandBuffer, 1)
    res := vk.AllocateCommandBuffers(
        a.device,
        &allocInfo,
        commandBuffers,
    )
    if res != vk.Success {
        return nil, fmt.Errorf("failed to allocate command buffer: %w", vk.Error(res))
    }
    commandBuffer := commandBuffers[0]

    beginInfo := vk.CommandBufferBeginInfo{
        SType: vk.StructureTypeCommandBufferBeginInfo,
        Flags: vk.CommandBufferUsageFlags(vk.CommandBufferUsageOneTimeSubmitBit),
    }

    vk.BeginCommandBuffer(commandBuffer, &beginInfo)

    return commandBuffer, nil
}

func (a *VulkanApp) endSingleTimeCommands(commandBuffer vk.CommandBuffer) error {
    commandBuffers := []vk.CommandBuffer{commandBuffer}

    defer func() {
        vk.FreeCommandBuffers(a.device, a.commandPool, 1, commandBuffers)
    }()

    res := vk.EndCommandBuffer(commandBuffer)
    if res != vk.Success {
        return fmt.Errorf("failed end command buffer: %w", vk.Error(res))
    }

    submitInfo := vk.SubmitInfo{
        SType:              vk.StructureTypeSubmitInfo,
        CommandBufferCount: 1,
        PCommandBuffers:    commandBuffers,
    }

    res = vk.QueueSubmit(a.graphicsQueue, 1, []vk.SubmitInfo{submitInfo}, vk.NullFence)
    if res != vk.Success {
        return fmt.Errorf("failed to submit to graphics queue: %w", vk.Error(res))
    }

    res = vk.QueueWaitIdle(a.graphicsQueue)
    if res != vk.Success {
        return fmt.Errorf("failed to wait on graphics queue idle: %w", vk.Error(res))
    }

    return nil
}

func (a *VulkanApp) transitionImageLayout(
    image vk.Image,
    oldLayout vk.ImageLayout,
    newLayout vk.ImageLayout,
) error {
    commandBuffer, err := a.beginSingleTimeCommands()
    if err != nil {
        return fmt.Errorf("failed to begin single time commands: %w", err)
    }

    barrier := vk.ImageMemoryBarrier{
        SType:               vk.StructureTypeImageMemoryBarrier,
        OldLayout:           oldLayout,
        NewLayout:           newLayout,
        SrcQueueFamilyIndex: vk.QueueFamilyIgnored,
        DstQueueFamilyIndex: vk.QueueFamilyIgnored,
        Image:               image,
        SubresourceRange: vk.ImageSubresourceRange{
            AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
            BaseMipLevel:   0,
            LevelCount:     1,
            BaseArrayLayer: 0,
            LayerCount:     1,
        },
        SrcAccessMask: 0,
        DstAccessMask: 0,
    }

    var (
        sourceStage      vk.PipelineStageFlags
        destinationStage vk.PipelineStageFlags
    )

    if oldLayout == vk.ImageLayoutUndefined &&
        newLayout == vk.ImageLayoutTransferDstOptimal {

        barrier.SrcAccessMask = 0
        barrier.DstAccessMask = vk.AccessFlags(vk.AccessHostWriteBit)

        sourceStage = vk.PipelineStageFlags(vk.PipelineStageTopOfPipeBit)
        destinationStage = vk.PipelineStageFlags(vk.PipelineStageHostBit)

    } else if oldLayout == vk.ImageLayoutUndefined &&
        newLayout == vk.ImageLayoutTransferSrcOptimal {

        barrier.SrcAccessMask = 0
        barrier.DstAccessMask = vk.AccessFlags(vk.AccessTransferReadBit)

        sourceStage = vk.PipelineStageFlags(vk.PipelineStageTopOfPipeBit)
        destinationStage = vk.PipelineStageFlags(vk.PipelineStageTransferBit)

    } else if oldLayout == vk.ImageLayoutTransferDstOptimal &&
        newLayout == vk.ImageLayoutShaderReadOnlyOptimal {

        barrier.SrcAccessMask = vk.AccessFlags(vk.AccessTransferWriteBit)
        barrier.DstAccessMask = vk.AccessFlags(vk.AccessShaderReadBit)

        sourceStage = vk.PipelineStageFlags(vk.PipelineStageTransferBit)
        destinationStage = vk.PipelineStageFlags(vk.PipelineStageFragmentShaderBit)

    } else if oldLayout == vk.ImageLayoutUndefined &&
        newLayout == vk.ImageLayoutShaderReadOnlyOptimal {

        barrier.SrcAccessMask = 0
        barrier.DstAccessMask = vk.AccessFlags(vk.AccessShaderReadBit)

        sourceStage = vk.PipelineStageFlags(vk.PipelineStageTopOfPipeBit)
        destinationStage = vk.PipelineStageFlags(vk.PipelineStageFragmentShaderBit)

    } else if oldLayout == vk.ImageLayoutTransferDstOptimal &&
        newLayout == vk.ImageLayoutTransferSrcOptimal {

        barrier.SrcAccessMask = vk.AccessFlags(vk.AccessHostWriteBit)
        barrier.DstAccessMask = vk.AccessFlags(vk.AccessTransferReadBit)

        sourceStage = vk.PipelineStageFlags(vk.PipelineStageHostBit)
        destinationStage = vk.PipelineStageFlags(vk.PipelineStageTransferBit)

    } else if oldLayout == vk.ImageLayoutTransferSrcOptimal &&
        newLayout == vk.ImageLayoutTransferDstOptimal {

        barrier.SrcAccessMask = vk.AccessFlags(vk.AccessTransferReadBit)
        barrier.DstAccessMask = vk.AccessFlags(vk.AccessHostWriteBit)

        sourceStage = vk.PipelineStageFlags(vk.PipelineStageTransferBit)
        destinationStage = vk.PipelineStageFlags(vk.PipelineStageHostBit)

    } else {
        return fmt.Errorf("unsupported layout transition from %d to %d",
            oldLayout, newLayout)
    }

    vk.CmdPipelineBarrier(
        commandBuffer,
        sourceStage, destinationStage,
        0,
        0, nil,
        0, nil,
        1, []vk.ImageMemoryBarrier{barrier},
    )

    return a.endSingleTimeCommands(commandBuffer)
}

func (a *VulkanApp) copyBufferToImage(
    buffer vk.Buffer,
    image vk.Image,
    width, height uint32,
) error {
    commandBuffer, err := a.beginSingleTimeCommands()
    if err != nil {
        return fmt.Errorf("failed to beging single time command buffer: %w", err)
    }

    region := vk.BufferImageCopy{
        BufferOffset:      0,
        BufferRowLength:   0,
        BufferImageHeight: 0,

        ImageSubresource: vk.ImageSubresourceLayers{
            AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
            MipLevel:       0,
            BaseArrayLayer: 0,
            LayerCount:     1,
        },

        ImageOffset: vk.Offset3D{
            X: 0, Y: 0, Z: 0,
        },

        ImageExtent: vk.Extent3D{
            Width:  width,
            Height: height,
            Depth:  1,
        },
    }

    vk.CmdCopyBufferToImage(
        commandBuffer,
        buffer,
        image,
        vk.ImageLayoutTransferDstOptimal,
        1,
        []vk.BufferImageCopy{region},
    )

    return a.endSingleTimeCommands(commandBuffer)
}

func (a *VulkanApp) createBuffer(
    size vk.DeviceSize,
    usage vk.BufferUsageFlags,
    properties vk.MemoryPropertyFlags,
    buffer *vk.Buffer,
    bufferMemory *vk.DeviceMemory,
) error {
    bufferInfo := vk.BufferCreateInfo{
        SType:       vk.StructureTypeBufferCreateInfo,
        Size:        size,
        Usage:       usage,
        SharingMode: vk.SharingModeExclusive,
    }

    res := vk.CreateBuffer(a.device, &bufferInfo, nil, buffer)
    if res != vk.Success {
        return fmt.Errorf("failed to create vertex buffer: %w", vk.Error(res))
    }

    var memRequirements vk.MemoryRequirements
    vk.GetBufferMemoryRequirements(a.device, *buffer, &memRequirements)
    memRequirements.Deref()

    memTypeIndex, err := a.findMemoryType(memRequirements.MemoryTypeBits, properties)
    if err != nil {
        return err
    }

    allocInfo := vk.MemoryAllocateInfo{
        SType:           vk.StructureTypeMemoryAllocateInfo,
        AllocationSize:  memRequirements.Size,
        MemoryTypeIndex: memTypeIndex,
    }

    res = vk.AllocateMemory(a.device, &allocInfo, nil, bufferMemory)
    if res != vk.Success {
        return fmt.Errorf("failed to allocate vertex buffer memory: %s", vk.Error(res))
    }

    res = vk.BindBufferMemory(a.device, *buffer, *bufferMemory, 0)
    if res != vk.Success {
        return fmt.Errorf("failed to bind buffer memory: %w", vk.Error(res))
    }

    return nil
}

func (a *VulkanApp) createImage(
    width uint32,
    height uint32,
    format vk.Format,
    tiling vk.ImageTiling,
    usage vk.ImageUsageFlags,
    properties vk.MemoryPropertyFlags,
    image *vk.Image,
    imageMemory *vk.DeviceMemory,
) error {
    imageInfo := vk.ImageCreateInfo{
        SType:     vk.StructureTypeImageCreateInfo,
        ImageType: vk.ImageType2d,
        Extent: vk.Extent3D{
            Width:  width,
            Height: height,
            Depth:  1,
        },
        MipLevels:     1,
        ArrayLayers:   1,
        Format:        format,
        Tiling:        tiling,
        InitialLayout: vk.ImageLayoutUndefined,
        Usage:         usage,
        SharingMode:   vk.SharingModeExclusive,
        Samples:       vk.SampleCount1Bit,
    }

    res := vk.CreateImage(a.device, &imageInfo, nil, image)
    if res != vk.Success {
        return fmt.Errorf("failed to create an image: %w", vk.Error(res))
    }

    var memRequirements vk.MemoryRequirements
    vk.GetImageMemoryRequirements(a.device, *image, &memRequirements)
    memRequirements.Deref()

    memTypeIndex, err := a.findMemoryType(memRequirements.MemoryTypeBits, properties)
    if err != nil {
        return err
    }

    allocInfo := vk.MemoryAllocateInfo{
        SType:           vk.StructureTypeMemoryAllocateInfo,
        AllocationSize:  memRequirements.Size,
        MemoryTypeIndex: memTypeIndex,
    }

    res = vk.AllocateMemory(a.device, &allocInfo, nil, imageMemory)
    if res != vk.Success {
        return fmt.Errorf("failed to allocate image buffer memory: %s", vk.Error(res))
    }

    res = vk.BindImageMemory(a.device, *image, *imageMemory, 0)
    if res != vk.Success {
        return fmt.Errorf("failed to bind image memory: %w", vk.Error(res))
    }

    return nil
}

func (a *VulkanApp) findMemoryType(
    typeFilter uint32,
    properties vk.MemoryPropertyFlags,
) (uint32, error) {
    var memProperties vk.PhysicalDeviceMemoryProperties
    vk.GetPhysicalDeviceMemoryProperties(a.physicalDevice, &memProperties)
    memProperties.Deref()

    for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
        memType := memProperties.MemoryTypes[i]
        memType.Deref()

        if typeFilter&(1<<i) == 0 {
            continue
        }

        if memType.PropertyFlags&properties != properties {
            continue
        }

        return i, nil
    }

    return 0, fmt.Errorf("failed to find suitable memory type")
}

func (a *VulkanApp) createFilmImage() error {
    texWidth := uint32(a.swapChainExtend.Width)
    texHeight := uint32(a.swapChainExtend.Height)

    a.film = newVulkanFilm(texWidth, texHeight)

    imgSize := vk.DeviceSize(a.film.getBufferSize())
    a.filmImageFormat = a.film.getFormat()

    var (
        stagingBuffer       vk.Buffer
        stagingBufferMemory vk.DeviceMemory
    )

    err := a.createBuffer(
        imgSize,
        vk.BufferUsageFlags(vk.BufferUsageTransferSrcBit),
        vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit)|
            vk.MemoryPropertyFlags(vk.MemoryPropertyHostCoherentBit),
        &stagingBuffer,
        &stagingBufferMemory,
    )
    if err != nil {
        return fmt.Errorf("failed to create staging GPU buffer: %w", err)
    }

    a.filmStagingBuffer = stagingBuffer
    a.filmStagingBufferMemory = stagingBufferMemory

    var pData unsafe.Pointer
    vk.MapMemory(a.device, a.filmStagingBufferMemory, 0, imgSize, 0, &pData)
    a.filmBufferData = pData

    var (
        filmImage       vk.Image
        filmImageMemory vk.DeviceMemory
    )

    err = a.createImage(
        texWidth,
        texHeight,
        a.filmImageFormat,
        vk.ImageTilingOptimal,
        vk.ImageUsageFlags(vk.ImageUsageTransferDstBit)|
            vk.ImageUsageFlags(vk.ImageUsageTransferSrcBit)|
            vk.ImageUsageFlags(vk.ImageUsageSampledBit),
        vk.MemoryPropertyFlags(vk.MemoryPropertyDeviceLocalBit),
        &filmImage,
        &filmImageMemory,
    )
    if err != nil {
        return fmt.Errorf("filed to create Vulkan image: %w", err)
    }
    a.filmImage = filmImage
    a.filmImageMemory = filmImageMemory

    err = a.transitionImageLayout(
        a.filmImage,
        vk.ImageLayoutUndefined,
        vk.ImageLayoutTransferSrcOptimal,
    )
    if err != nil {
        return fmt.Errorf("transition image layout: %w", err)
    }

    return nil
}

func (a *VulkanApp) cleanFilmImage() {
    if a.filmImage != vk.NullImage {
        vk.DestroyImage(a.device, a.filmImage, nil)
    }

    if a.filmImageMemory != vk.NullDeviceMemory {
        vk.FreeMemory(a.device, a.filmImageMemory, nil)
    }

    if a.filmStagingBuffer != vk.NullBuffer {
        vk.DestroyBuffer(a.device, a.filmStagingBuffer, nil)
    }

    if a.filmStagingBufferMemory != vk.NullDeviceMemory {
        vk.UnmapMemory(a.device, a.filmStagingBufferMemory)
        vk.FreeMemory(a.device, a.filmStagingBufferMemory, nil)
    }
}

func (a *VulkanApp) copyFilmToGPUImage() error {
    filmBytes := a.film.asVkBuffer()

    // copy data to the staging buffer
    vk.Memcopy(a.filmBufferData, filmBytes)

    width, height := a.swapChainExtend.Width, a.swapChainExtend.Height

    err := a.transitionImageLayout(
        a.filmImage,
        vk.ImageLayoutTransferSrcOptimal,
        vk.ImageLayoutTransferDstOptimal,
    )
    if err != nil {
        return fmt.Errorf("transition image layout (dst): %w", err)
    }

    err = a.copyBufferToImage(a.filmStagingBuffer, a.filmImage, width, height)
    if err != nil {
        return fmt.Errorf("copying buffer to image: %w", err)
    }

    err = a.transitionImageLayout(
        a.filmImage,
        vk.ImageLayoutTransferDstOptimal,
        vk.ImageLayoutTransferSrcOptimal,
    )
    if err != nil {
        return fmt.Errorf("transition image layout (src): %w", err)
    }

    return nil
}

func (a *VulkanApp) createCommandPool() error {
    queueFamilyIndices := a.findQueueFamilies(a.physicalDevice)
    poolInfo := vk.CommandPoolCreateInfo{
        SType: vk.StructureTypeCommandPoolCreateInfo,
        Flags: vk.CommandPoolCreateFlags(
            vk.CommandPoolCreateResetCommandBufferBit,
        ),
        QueueFamilyIndex: queueFamilyIndices.Graphics.Get(),
    }

    var commandPool vk.CommandPool
    res := vk.CreateCommandPool(a.device, &poolInfo, nil, &commandPool)
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("failed to create command pool: %w", err)
    }
    a.commandPool = commandPool

    return nil
}

func (a *VulkanApp) createCommandBuffer() error {
    allocInfo := vk.CommandBufferAllocateInfo{
        SType:              vk.StructureTypeCommandBufferAllocateInfo,
        CommandPool:        a.commandPool,
        Level:              vk.CommandBufferLevelPrimary,
        CommandBufferCount: 2,
    }

    commandBuffers := make([]vk.CommandBuffer, 2)
    res := vk.AllocateCommandBuffers(a.device, &allocInfo, commandBuffers)
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("failed to allocate command buffer: %w", err)
    }
    a.commandBuffers = commandBuffers

    return nil
}

func (a *VulkanApp) recordCommandBuffer(
    commandBuffer vk.CommandBuffer,
    imageIndex uint32,
) error {
    beginInfo := vk.CommandBufferBeginInfo{
        SType: vk.StructureTypeCommandBufferBeginInfo,
        Flags: 0,
    }

    res := vk.BeginCommandBuffer(commandBuffer, &beginInfo)
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("cannot add begin command to the buffer: %w", err)
    }

    clearColor := vk.NewClearValue([]float32{0, 0, 0, 1})

    renderPassInfo := vk.RenderPassBeginInfo{
        SType:       vk.StructureTypeRenderPassBeginInfo,
        RenderPass:  a.renderPass,
        Framebuffer: a.swapChainFramebuffers[imageIndex],
        RenderArea: vk.Rect2D{
            Offset: vk.Offset2D{
                X: 0,
                Y: 0,
            },
            Extent: a.swapChainExtend,
        },
        ClearValueCount: 1,
        PClearValues:    []vk.ClearValue{clearColor},
    }

    vk.CmdBeginRenderPass(commandBuffer, &renderPassInfo, vk.SubpassContentsInline)
    vk.CmdBindPipeline(commandBuffer, vk.PipelineBindPointGraphics, a.graphicsPipline)

    viewport := vk.Viewport{
        X: 0, Y: 0,
        Width:    float32(a.swapChainExtend.Width),
        Height:   float32(a.swapChainExtend.Height),
        MinDepth: 0,
        MaxDepth: 1,
    }
    vk.CmdSetViewport(commandBuffer, 0, 1, []vk.Viewport{viewport})

    scissor := vk.Rect2D{
        Offset: vk.Offset2D{X: 0, Y: 0},
        Extent: a.swapChainExtend,
    }
    vk.CmdSetScissor(commandBuffer, 0, 1, []vk.Rect2D{scissor})

    vk.CmdEndRenderPass(commandBuffer)

    barrier := vk.ImageMemoryBarrier{
        SType:               vk.StructureTypeImageMemoryBarrier,
        OldLayout:           vk.ImageLayoutPresentSrc,
        NewLayout:           vk.ImageLayoutTransferDstOptimal,
        SrcQueueFamilyIndex: vk.QueueFamilyIgnored,
        DstQueueFamilyIndex: vk.QueueFamilyIgnored,
        Image:               a.swapChainImages[imageIndex],
        SubresourceRange: vk.ImageSubresourceRange{
            AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
            BaseMipLevel:   0,
            LevelCount:     1,
            BaseArrayLayer: 0,
            LayerCount:     1,
        },
        SrcAccessMask: vk.AccessFlags(vk.AccessShaderWriteBit),
        DstAccessMask: vk.AccessFlags(vk.AccessTransferWriteBit),
    }

    vk.CmdPipelineBarrier(
        commandBuffer,
        vk.PipelineStageFlags(vk.PipelineStageVertexShaderBit),
        vk.PipelineStageFlags(vk.PipelineStageTransferBit),
        0,
        0, nil,
        0, nil,
        1, []vk.ImageMemoryBarrier{barrier},
    )

    blitRegion := vk.ImageBlit{
        SrcSubresource: vk.ImageSubresourceLayers{
            AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
            MipLevel:       0,
            BaseArrayLayer: 0,
            LayerCount:     1,
        },
        SrcOffsets: [2]vk.Offset3D{
            {X: 0, Y: 0, Z: 0},
            {X: int32(a.swapChainExtend.Width), Y: int32(a.swapChainExtend.Height), Z: 1},
        },

        DstSubresource: vk.ImageSubresourceLayers{
            AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
            MipLevel:       0,
            BaseArrayLayer: 0,
            LayerCount:     1,
        },
        DstOffsets: [2]vk.Offset3D{
            {X: 0, Y: 0, Z: 0},
            {X: int32(a.swapChainExtend.Width), Y: int32(a.swapChainExtend.Height), Z: 1},
        },
    }

    vk.CmdBlitImage(
        commandBuffer,
        a.filmImage,
        vk.ImageLayoutTransferSrcOptimal,
        a.swapChainImages[imageIndex],
        vk.ImageLayoutTransferDstOptimal,
        1,
        []vk.ImageBlit{blitRegion},
        vk.FilterLinear,
    )

    barrier = vk.ImageMemoryBarrier{
        SType:               vk.StructureTypeImageMemoryBarrier,
        OldLayout:           vk.ImageLayoutTransferDstOptimal,
        NewLayout:           vk.ImageLayoutPresentSrc,
        SrcQueueFamilyIndex: vk.QueueFamilyIgnored,
        DstQueueFamilyIndex: vk.QueueFamilyIgnored,
        Image:               a.swapChainImages[imageIndex],
        SubresourceRange: vk.ImageSubresourceRange{
            AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
            BaseMipLevel:   0,
            LevelCount:     1,
            BaseArrayLayer: 0,
            LayerCount:     1,
        },
        SrcAccessMask: vk.AccessFlags(vk.AccessTransferWriteBit),
        DstAccessMask: vk.AccessFlags(vk.AccessShaderWriteBit),
    }

    vk.CmdPipelineBarrier(
        commandBuffer,
        vk.PipelineStageFlags(vk.PipelineStageTransferBit),
        vk.PipelineStageFlags(vk.PipelineStageVertexShaderBit),
        0,
        0, nil,
        0, nil,
        1, []vk.ImageMemoryBarrier{barrier},
    )

    if err := vk.Error(vk.EndCommandBuffer(commandBuffer)); err != nil {
        return fmt.Errorf("recording commands to buffer failed: %w", err)
    }
    return nil
}

func (a *VulkanApp) createSyncObjects() error {
    semaphoreInfo := vk.SemaphoreCreateInfo{
        SType: vk.StructureTypeSemaphoreCreateInfo,
    }

    fenceInfo := vk.FenceCreateInfo{
        SType: vk.StructureTypeFenceCreateInfo,
        Flags: vk.FenceCreateFlags(vk.FenceCreateSignaledBit),
    }

    for i := 0; i < maxFramesInFlight; i++ {
        var imageAvailabmeSem vk.Semaphore
        if err := vk.Error(
            vk.CreateSemaphore(a.device, &semaphoreInfo, nil, &imageAvailabmeSem),
        ); err != nil {
            return fmt.Errorf("failed to create imageAvailabmeSem: %w", err)
        }
        a.imageAvailabmeSems = append(a.imageAvailabmeSems, imageAvailabmeSem)

        var renderFinishedSem vk.Semaphore
        if err := vk.Error(
            vk.CreateSemaphore(a.device, &semaphoreInfo, nil, &renderFinishedSem),
        ); err != nil {
            return fmt.Errorf("failed to create renderFinishedSem: %w", err)
        }
        a.renderFinishedSems = append(a.renderFinishedSems, renderFinishedSem)

        var fence vk.Fence
        if err := vk.Error(
            vk.CreateFence(a.device, &fenceInfo, nil, &fence),
        ); err != nil {
            return fmt.Errorf("failed to create inFlightFence: %w", err)
        }
        a.inFlightFences = append(a.inFlightFences, fence)
    }

    return nil
}

func (a *VulkanApp) createInstance() error {
    if a.enableValidationLayers && !a.checkValidationSupport() {
        return fmt.Errorf("validation layers requested but not available")
    }

    appInfo := vk.ApplicationInfo{
        SType:              vk.StructureTypeApplicationInfo,
        PApplicationName:   title + "\x00",
        ApplicationVersion: vk.MakeVersion(1, 0, 0),
        PEngineName:        "No Engine\x00",
        EngineVersion:      vk.MakeVersion(1, 0, 0),
        ApiVersion:         vk.ApiVersion10,
    }

    glfwExtensions := glfw.GetCurrentContext().GetRequiredInstanceExtensions()
    createInfo := vk.InstanceCreateInfo{
        SType:                   vk.StructureTypeInstanceCreateInfo,
        PApplicationInfo:        &appInfo,
        EnabledExtensionCount:   uint32(len(glfwExtensions)),
        PpEnabledExtensionNames: glfwExtensions,
    }

    if a.enableValidationLayers {
        createInfo.EnabledLayerCount = uint32(len(a.validationLayers))
        createInfo.PpEnabledLayerNames = a.validationLayers
    }

    var instance vk.Instance
    if res := vk.CreateInstance(&createInfo, nil, &instance); res != vk.Success {
        return fmt.Errorf("failed to create Vulkan instance: %w", vk.Error(res))
    }

    a.instance = instance
    return nil
}

// findQueueFamilies returns a FamilyIndeces populated with Vulkan queue families needed
// by the program.
func (a *VulkanApp) findQueueFamilies(
    device vk.PhysicalDevice,
) QueueFamilyIndices {
    indices := QueueFamilyIndices{}

    var queueFamilyCount uint32
    vk.GetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, nil)

    queueFamilies := make([]vk.QueueFamilyProperties, queueFamilyCount)
    vk.GetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, queueFamilies)

    for i, family := range queueFamilies {
        family.Deref()

        if family.QueueFlags&vk.QueueFlags(vk.QueueGraphicsBit) != 0 {
            indices.Graphics.Set(uint32(i))
        }

        var hasPresent vk.Bool32
        err := vk.Error(
            vk.GetPhysicalDeviceSurfaceSupport(device, uint32(i), a.surface, &hasPresent),
        )
        if err != nil {
            fmt.Printf("error querying surface support for queue family %d: %s\n", i, err)
        } else if hasPresent.B() {
            indices.Present.Set(uint32(i))
        }

        if indices.IsComplete() {
            break
        }
    }

    return indices
}

func (a *VulkanApp) querySwapChainSupport(
    device vk.PhysicalDevice,
) swapChainSupportDetails {
    details := swapChainSupportDetails{}

    var capabilities vk.SurfaceCapabilities
    res := vk.GetPhysicalDeviceSurfaceCapabilities(device, a.surface, &capabilities)
    if err := vk.Error(res); err != nil {
        panic(fmt.Sprintf("failed to query device surface capabilities: %s", err))
    }
    capabilities.Deref()
    capabilities.CurrentExtent.Deref()
    capabilities.MinImageExtent.Deref()
    capabilities.MaxImageExtent.Deref()

    details.capabilities = capabilities

    var formatCount uint32
    res = vk.GetPhysicalDeviceSurfaceFormats(device, a.surface, &formatCount, nil)
    if err := vk.Error(res); err != nil {
        panic(fmt.Sprintf("failed to query device surface formats: %s", err))
    }

    if formatCount != 0 {
        formats := make([]vk.SurfaceFormat, formatCount)
        vk.GetPhysicalDeviceSurfaceFormats(device, a.surface, &formatCount, formats)
        for _, format := range formats {
            format.Deref()
            details.formats = append(details.formats, format)
        }
    }

    var presentModeCount uint32
    res = vk.GetPhysicalDeviceSurfacePresentModes(
        device, a.surface, &presentModeCount, nil,
    )
    if err := vk.Error(res); err != nil {
        panic(fmt.Sprintf("failed to query device surface present modes: %s", err))
    }

    if presentModeCount != 0 {
        presentModes := make([]vk.PresentMode, presentModeCount)
        vk.GetPhysicalDeviceSurfacePresentModes(
            device, a.surface, &presentModeCount, presentModes,
        )
        details.presentModes = presentModes
    }

    return details
}

// getDeviceScore returns how suitable is this device for the current program.
// Bigger score means better. Zero or negative means the device cannot be used.
func (a *VulkanApp) getDeviceScore(device vk.PhysicalDevice) uint32 {
    var (
        deviceScore uint32
        properties  vk.PhysicalDeviceProperties
    )

    vk.GetPhysicalDeviceProperties(device, &properties)
    properties.Deref()

    if properties.DeviceType == vk.PhysicalDeviceTypeDiscreteGpu {
        deviceScore += 1000
    } else {
        deviceScore++
    }

    if !a.isDeviceSuitable(device) {
        deviceScore = 0
    }

    if a.args.Debug {
        fmt.Printf(
            "Available device: %s (score: %d)\n",
            vk.ToString(properties.DeviceName[:]),
            deviceScore,
        )
    }

    return deviceScore
}

func (a *VulkanApp) isDeviceSuitable(device vk.PhysicalDevice) bool {
    indices := a.findQueueFamilies(device)
    extensionsSupported := a.checkDeviceExtensionSupport(device)

    swapChainAdequate := false
    if extensionsSupported {
        swapChainSupport := a.querySwapChainSupport(device)
        swapChainAdequate = len(swapChainSupport.formats) > 0 &&
            len(swapChainSupport.presentModes) > 0
    }

    return indices.IsComplete() && extensionsSupported && swapChainAdequate
}

func (a *VulkanApp) chooseSwapSurfaceFormat(
    availableFormats []vk.SurfaceFormat,
) vk.SurfaceFormat {
    if a.args.Debug {
        fmt.Println("Available swapchanin formats:")
        for _, format := range availableFormats {
            fmt.Printf("\t* %#v\n", format.Format)
        }
    }

    for _, format := range availableFormats {
        if format.Format == vk.FormatB8g8r8a8Srgb &&
            format.ColorSpace == vk.ColorSpaceSrgbNonlinear {
            return format
        }
    }

    return availableFormats[0]
}

func (a *VulkanApp) chooseSwapPresentMode(
    available []vk.PresentMode,
) vk.PresentMode {
    if a.args.VSync {
        return vk.PresentModeFifo
    }

    for _, mode := range available {
        if mode == vk.PresentModeMailbox {
            return mode
        }
    }

    return vk.PresentModeFifo
}

func (a *VulkanApp) chooseSwapExtend(
    capabilities vk.SurfaceCapabilities,
) vk.Extent2D {
    if capabilities.CurrentExtent.Width != math.MaxUint32 {
        return capabilities.CurrentExtent
    }

    width, height := a.window.GetFramebufferSize()

    actualExtend := vk.Extent2D{
        Width:  uint32(width),
        Height: uint32(height),
    }

    actualExtend.Width = clamp(
        actualExtend.Width,
        capabilities.MinImageExtent.Width,
        capabilities.MaxImageExtent.Width,
    )

    actualExtend.Height = clamp(
        actualExtend.Height,
        capabilities.MinImageExtent.Height,
        capabilities.MaxImageExtent.Height,
    )

    fmt.Printf("actualExtend: %#v", actualExtend)

    return actualExtend
}

func (a *VulkanApp) checkDeviceExtensionSupport(device vk.PhysicalDevice) bool {
    var extensionsCount uint32
    res := vk.EnumerateDeviceExtensionProperties(device, "", &extensionsCount, nil)
    if err := vk.Error(res); err != nil {
        fmt.Printf(
            "WARNING: enumerating device (%d) extension properties count: %s\n",
            device,
            err,
        )
        return false
    }

    availableExtensions := make([]vk.ExtensionProperties, extensionsCount)
    res = vk.EnumerateDeviceExtensionProperties(device, "", &extensionsCount,
        availableExtensions)
    if err := vk.Error(res); err != nil {
        fmt.Printf("WARNING: getting device (%d) extension properties: %s\n", device, err)
        return false
    }

    requiredExtensions := make(map[string]struct{})
    for _, extensionName := range a.deviceExtensions {
        requiredExtensions[extensionName] = struct{}{}
    }

    for _, extension := range availableExtensions {
        extension.Deref()
        extensionName := vk.ToString(extension.ExtensionName[:])

        delete(requiredExtensions, extensionName+"\x00")
    }

    return len(requiredExtensions) == 0
}

func (a *VulkanApp) checkValidationSupport() bool {
    var count uint32
    if vk.EnumerateInstanceLayerProperties(&count, nil) != vk.Success {
        return false
    }
    availableLayers := make([]vk.LayerProperties, count)

    if vk.EnumerateInstanceLayerProperties(&count, availableLayers) != vk.Success {
        return false
    }

    availableLayersStr := make([]string, 0, count)
    for _, layer := range availableLayers {
        layer.Deref()

        layerName := vk.ToString(layer.LayerName[:])
        availableLayersStr = append(availableLayersStr, layerName+"\x00")
    }

    for _, validationLayer := range a.validationLayers {
        if !slices.Contains(availableLayersStr, validationLayer) {
            return false
        }
    }

    return true
}

func (a *VulkanApp) createShaderModule(code []byte) (vk.ShaderModule, error) {
    createInfo := vk.ShaderModuleCreateInfo{
        SType:    vk.StructureTypeShaderModuleCreateInfo,
        CodeSize: uint(len(code)),
        PCode:    unsafer.SliceBytesToUint32(code),
    }

    var shaderModule vk.ShaderModule
    res := vk.CreateShaderModule(a.device, &createInfo, nil, &shaderModule)
    return shaderModule, vk.Error(res)
}

func (a *VulkanApp) initEngine() error {
    width, height := a.swapChainExtend.Width, a.swapChainExtend.Height

    smpl := sampler.NewSimple(int(width), int(height), a.film)

    if a.args.Interactive {
        smpl.MakeContinuous()
    }

    cam := scene.GetCamera(float64(width), float64(height))

    tracer := engine.NewFPS(smpl)
    tracer.SetTarget(a.film, cam)
    tracer.ShowBBoxes = a.args.ShowBBoxes

    fmt.Printf("Loading scene...\n")
    loadingStart := time.Now()
    tracer.Scene.InitScene(a.args.SceneName)
    fmt.Printf("Loading scene took %s\n", time.Since(loadingStart))

    a.sampler = smpl
    a.tracer = tracer
    a.cam = cam

    return nil
}

func (a *VulkanApp) cleanEngine() {
    a.sampler.Stop()
    a.tracer.StopRendering()
}

func (a *VulkanApp) mainLoop() error {
    a.tracer.Render()
    minFrameTime := time.Duration(1000.0/float32(a.args.FPSCap)) * time.Millisecond

    var (
        traceStarted bool
        bPressed     bool

        frameCounter uint64
        lastShowFPS  = time.Now()

        // When `dirty` is "false" after a full raytraced frame is completed then tracing
        // could stop for a bit and wit for some movement before continuing.
        dirty     bool = true
        prevDrity bool
    )

    for !a.window.ShouldClose() {
        renderStart := time.Now()
        err := a.drawFrame()
        if err != nil {
            return fmt.Errorf("error drawing a frame: %w", err)
        }
        renderTime := time.Since(renderStart)

        glfw.PollEvents()
        if a.args.Interactive {
            if handleInteractionEvents(a.window, a.cam, renderTime) {
                dirty = true
            }

            if !bPressed && a.window.GetKey(glfw.KeyB) == glfw.Press {
                a.tracer.ShowBBoxes = !a.tracer.ShowBBoxes
                bPressed = true
                dirty = true
            }

            if bPressed && a.window.GetKey(glfw.KeyB) == glfw.Release {
                bPressed = false
                dirty = true
            }

            if !traceStarted && a.window.GetKey(glfw.KeyT) == glfw.Press {
                dirty = true
                traceStarted = true
                go func() {
                    collectTrace()
                    traceStarted = false
                }()
            }
        }

        elapsed := time.Since(renderStart)

        if !a.args.VSync && a.args.FPSCap > 0 && elapsed < minFrameTime {
            time.Sleep(minFrameTime - elapsed)
        }

        if a.args.ShowFPS {
            frameCounter++
            now := time.Now()
            elapsed = now.Sub(lastShowFPS)

            if elapsed > time.Second {
                fps := float64(frameCounter) / elapsed.Seconds()
                fmt.Printf("\r                                                               ")
                fmt.Printf("\rFPS: %5.3f Render time: %8s Last frame: %12s",
                    fps, renderTime, a.film.FrameTime(),
                )

                frameCounter = 0
                lastShowFPS = time.Now()
            }
        }

        if dirty {
            if !prevDrity {
                a.tracer.Resume()
            }
            prevDrity = dirty
            dirty = false
        } else {
            if prevDrity {
                a.tracer.Pause()
            }
            prevDrity = dirty
        }
    }

    a.tracer.Resume()
    fmt.Println("\nClosing window, rendering stopped.")
    vk.DeviceWaitIdle(a.device)
    return nil
}

func (a *VulkanApp) drawFrame() error {
    if err := a.copyFilmToGPUImage(); err != nil {
        return fmt.Errorf("error copying film to GPU image: %w", err)
    }

    fences := []vk.Fence{a.inFlightFences[a.curentFrame]}
    vk.WaitForFences(a.device, 1, fences, vk.True, math.MaxUint64)

    var imageIndex uint32
    res := vk.AcquireNextImage(
        a.device,
        a.swapChain,
        math.MaxUint64,
        a.imageAvailabmeSems[a.curentFrame],
        vk.Fence(vk.NullHandle),
        &imageIndex,
    )
    if res == vk.ErrorOutOfDate {
        a.recreateSwapChain()
        return nil
    } else if res != vk.Success && res != vk.Suboptimal {
        return fmt.Errorf("failed to acquire swap chain image: %w", vk.Error(res))
    }

    // Only reset the fence if we are submitting work.
    vk.ResetFences(a.device, 1, fences)

    commandBuffer := a.commandBuffers[a.curentFrame]

    vk.ResetCommandBuffer(commandBuffer, 0)
    if err := a.recordCommandBuffer(commandBuffer, imageIndex); err != nil {
        return fmt.Errorf("recording command buffer: %w", err)
    }

    signalSemaphores := []vk.Semaphore{
        a.renderFinishedSems[a.curentFrame],
    }

    submitInfo := vk.SubmitInfo{
        SType:              vk.StructureTypeSubmitInfo,
        WaitSemaphoreCount: 1,
        PWaitSemaphores:    []vk.Semaphore{a.imageAvailabmeSems[a.curentFrame]},
        PWaitDstStageMask: []vk.PipelineStageFlags{
            vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
        },
        CommandBufferCount:   1,
        PCommandBuffers:      []vk.CommandBuffer{commandBuffer},
        PSignalSemaphores:    signalSemaphores,
        SignalSemaphoreCount: uint32(len(signalSemaphores)),
    }

    res = vk.QueueSubmit(
        a.graphicsQueue,
        1,
        []vk.SubmitInfo{submitInfo},
        a.inFlightFences[a.curentFrame],
    )
    if err := vk.Error(res); err != nil {
        return fmt.Errorf("queue submit error: %w", err)
    }

    swapChains := []vk.Swapchain{
        a.swapChain,
    }

    presentInfo := vk.PresentInfo{
        SType:              vk.StructureTypePresentInfo,
        WaitSemaphoreCount: uint32(len(signalSemaphores)),
        PWaitSemaphores:    signalSemaphores,
        SwapchainCount:     uint32(len(swapChains)),
        PSwapchains:        swapChains,
        PImageIndices:      []uint32{imageIndex},
    }

    res = vk.QueuePresent(a.presentQueue, &presentInfo)
    if res == vk.ErrorOutOfDate || res == vk.Suboptimal || a.frameBufferResized {
        a.frameBufferResized = false
        a.recreateSwapChain()
    } else if res != vk.Success {
        return fmt.Errorf("failed to present swap chain image: %w", vk.Error(res))
    }

    a.curentFrame = (a.curentFrame + 1) % maxFramesInFlight

    return nil
}

// swapChainSupportDetails describes a present surface. The type is suitable for
// passing around many details of the service between functions.
type swapChainSupportDetails struct {
    capabilities vk.SurfaceCapabilities
    formats      []vk.SurfaceFormat
    presentModes []vk.PresentMode
}

func clamp[T cmp.Ordered](val, min, max T) T {
    if val < min {
        val = min
    }
    if val > max {
        val = max
    }
    return val
}

// QueueFamilyIndices holds the indexes of Vulkan queue families needed by the programs.
type QueueFamilyIndices struct {

    // Graphics is the index of the graphics queue family.
    Graphics optional.Optional[uint32]

    // Present is the index of the queue family used for presenting to the drawing
    // surface.
    Present optional.Optional[uint32]
}

// IsComplete returns true if all families have been set.
func (f *QueueFamilyIndices) IsComplete() bool {
    return f.Graphics.HasValue() && f.Present.HasValue()
}
