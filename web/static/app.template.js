// Global variables
let wasmReady = false;
let selectedFiles = [];
let processedImages = [];
const WASM_VERSION = '{{VERSION}}'; // Injected at build time

// Initialize WebAssembly
async function initWasm() {
    const go = new Go();
    
    try {
        console.log('Loading WebAssembly...');
        // Add version query parameter for cache busting
        const response = await fetch(`static/imagekit.wasm?v=${WASM_VERSION}`);
        if (!response.ok) {
            throw new Error(`Failed to fetch WASM: ${response.status}`);
        }
        
        const buffer = await response.arrayBuffer();
        console.log('WASM file loaded, size:', buffer.byteLength);
        
        const result = await WebAssembly.instantiate(buffer, go.importObject);
        go.run(result.instance);
        
        // Wait a bit for Go to initialize
        await new Promise(resolve => setTimeout(resolve, 100));
        
        // Verify functions are available
        if (typeof window.processImage !== 'function') {
            throw new Error('WASM functions not registered');
        }
        
        wasmReady = true;
        document.getElementById('loading').classList.add('hidden');
        // Reset loading message for future use
        const initialLoadMessage = document.querySelector('.initial-load-message');
        if (initialLoadMessage) {
            initialLoadMessage.style.display = ''; // Reset display
        }
        
        // Update version display
        if (typeof window.getImageKitVersion === 'function') {
            const version = window.getImageKitVersion();
            const versionElement = document.getElementById('imagekit-version');
            if (versionElement) {
                versionElement.textContent = version;
            }
        }
        
        console.log('WebAssembly loaded successfully');
        console.log('Available functions:', Object.keys(window).filter(k => k.includes('Image')));
    } catch (error) {
        console.error('Failed to load WebAssembly:', error);
        document.getElementById('loading').innerHTML = `
            <p style="color: red;">WebAssembly 로드 실패</p>
            <p style="color: white; font-size: 0.9em;">${error.message}</p>
        `;
    }
}

// Register Service Worker for PWA
async function registerServiceWorker() {
    if ('serviceWorker' in navigator) {
        try {
            const registration = await navigator.serviceWorker.register('/static/service-worker.js');
            console.log('Service Worker registered:', registration);
            
            // Check for updates periodically
            setInterval(() => {
                registration.update();
            }, 60 * 60 * 1000); // Check every hour
            
            // Handle updates
            registration.addEventListener('updatefound', () => {
                const newWorker = registration.installing;
                newWorker.addEventListener('statechange', () => {
                    if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                        // New service worker available
                        if (confirm('새 버전이 있습니다. 업데이트하시겠습니까?')) {
                            newWorker.postMessage({ type: 'SKIP_WAITING' });
                            window.location.reload();
                        }
                    }
                });
            });
        } catch (error) {
            console.error('Service Worker registration failed:', error);
        }
    }
}

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initWasm();
    setupEventListeners();
    registerServiceWorker();
    
    // Set current year in footer
    const yearElement = document.getElementById('current-year');
    if (yearElement) {
        yearElement.textContent = new Date().getFullYear();
    }
    
    // Initialize crop dialog
    cropDialog = new CropDialog();
    
    // Check if app is running in standalone mode (PWA)
    const isInStandaloneMode = window.matchMedia('(display-mode: standalone)').matches || 
                               window.navigator.standalone || 
                               document.referrer.includes('android-app://');
    
    // Add install prompt handler for PWA
    let deferredPrompt;
    
    // Only set up install prompt if not already in PWA mode
    if (!isInStandaloneMode) {
        window.addEventListener('beforeinstallprompt', (e) => {
            e.preventDefault();
            deferredPrompt = e;
            
            // Get elements
            const container = document.getElementById('installPWAContainer');
            const fab = document.getElementById('installPWAFab');
            const expanded = document.getElementById('installPWAExpanded');
            const installBtn = document.getElementById('installBtn');
            const minimizeBtn = document.getElementById('minimizeBtn');
            
            if (container && fab && expanded) {
                // Show the FAB container
                container.style.display = 'block';
            
            // Check if user has previously minimized (not permanently dismissed)
            const wasMinimized = localStorage.getItem('pwa-install-minimized');
            
            // Show expanded view initially for first-time users
            if (!wasMinimized) {
                setTimeout(() => {
                    expanded.style.display = 'block';
                    // Auto-minimize after 8 seconds if not interacted
                    setTimeout(() => {
                        if (expanded.style.display === 'block') {
                            expanded.style.display = 'none';
                            localStorage.setItem('pwa-install-minimized', 'true');
                        }
                    }, 8000);
                }, 1000); // Show after 1 second delay
            }
            
            // FAB click handler - toggle expanded view
            fab.addEventListener('click', () => {
                if (expanded.style.display === 'none') {
                    expanded.style.display = 'block';
                    // Add scale animation
                    expanded.style.transform = 'scale(0.8)';
                    expanded.style.opacity = '0';
                    setTimeout(() => {
                        expanded.style.transform = 'scale(1)';
                        expanded.style.opacity = '1';
                    }, 10);
                } else {
                    expanded.style.transform = 'scale(0.8)';
                    expanded.style.opacity = '0';
                    setTimeout(() => {
                        expanded.style.display = 'none';
                        expanded.style.transform = 'scale(1)';
                        expanded.style.opacity = '1';
                    }, 300);
                }
            });
            
            // Install button handler
            if (installBtn) {
                installBtn.addEventListener('click', async () => {
                    if (deferredPrompt) {
                        deferredPrompt.prompt();
                        const { outcome } = await deferredPrompt.userChoice;
                        console.log('Install prompt outcome:', outcome);
                        
                        if (outcome === 'accepted') {
                            // Hide the FAB after successful installation
                            container.style.display = 'none';
                            localStorage.setItem('pwa-installed', 'true');
                        } else {
                            // User declined, just minimize
                            expanded.style.display = 'none';
                        }
                        
                        deferredPrompt = null;
                    }
                });
            }
            
            // Minimize button handler
            if (minimizeBtn) {
                minimizeBtn.addEventListener('click', () => {
                    expanded.style.transform = 'scale(0.8)';
                    expanded.style.opacity = '0';
                    setTimeout(() => {
                        expanded.style.display = 'none';
                        expanded.style.transform = 'scale(1)';
                        expanded.style.opacity = '1';
                    }, 300);
                    localStorage.setItem('pwa-install-minimized', 'true');
                });
            }
        }
    });
    } else {
        // App is already running in PWA mode
        console.log('App is running in standalone/PWA mode - install button hidden');
    }
});

// Setup event listeners
function setupEventListeners() {
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    
    // Load saved options from localStorage
    loadSavedOptions();
    
    // File upload
    if (uploadArea && fileInput) {
        uploadArea.addEventListener('click', () => fileInput.click());
        fileInput.addEventListener('change', handleFileSelect);
        
        // Drag and drop
        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadArea.classList.add('dragover');
        });
        
        uploadArea.addEventListener('dragleave', () => {
            uploadArea.classList.remove('dragover');
        });
        
        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('dragover');
            handleFiles(e.dataTransfer.files);
        });
    }
    
    // Options toggles with auto-save
    const enableResize = document.getElementById('enableResize');
    const resizeOptions = document.getElementById('resizeOptions');
    if (enableResize && resizeOptions) {
        enableResize.addEventListener('change', (e) => {
            if (e.target.checked) {
                resizeOptions.classList.remove('hidden');
            } else {
                resizeOptions.classList.add('hidden');
            }
            saveOptions();
        });
    }
    
    const enableDPI = document.getElementById('enableDPI');
    const dpiOptions = document.getElementById('dpiOptions');
    if (enableDPI && dpiOptions) {
        enableDPI.addEventListener('change', (e) => {
            if (e.target.checked) {
                dpiOptions.classList.remove('hidden');
            } else {
                dpiOptions.classList.add('hidden');
            }
            saveOptions();
        });
    }
    
    // Save options on input change
    const resizeWidth = document.getElementById('resizeWidth');
    const resizeHeight = document.getElementById('resizeHeight');
    const dpiValue = document.getElementById('dpiValue');
    
    if (resizeWidth) {
        resizeWidth.addEventListener('input', saveOptions);
    }
    if (resizeHeight) {
        resizeHeight.addEventListener('input', saveOptions);
    }
    if (dpiValue) {
        dpiValue.addEventListener('change', saveOptions);
    }
    
    // Process button
    const processBtn = document.getElementById('processBtn');
    if (processBtn) {
        processBtn.addEventListener('click', processImages);
    }
    
    // Download all button
    const downloadAllBtn = document.getElementById('downloadAllBtn');
    if (downloadAllBtn) {
        downloadAllBtn.addEventListener('click', downloadAll);
    }
}

// Load saved options from localStorage
function loadSavedOptions() {
    const saved = localStorage.getItem('imagekitOptions');
    
    if (saved) {
        try {
            const options = JSON.parse(saved);
            
            // Load resize options
            const enableResize = document.getElementById('enableResize');
            const resizeWidth = document.getElementById('resizeWidth');
            const resizeHeight = document.getElementById('resizeHeight');
            const resizeOptions = document.getElementById('resizeOptions');
            
            if (enableResize) {
                enableResize.checked = options.enableResize !== false; // Default true
            }
            if (resizeWidth) {
                resizeWidth.value = options.resizeWidth || '';
            }
            if (resizeHeight) {
                resizeHeight.value = options.resizeHeight || '';
            }
            if (resizeOptions) {
                if (options.enableResize !== false) {
                    resizeOptions.classList.remove('hidden');
                } else {
                    resizeOptions.classList.add('hidden');
                }
            }
            
            // Load DPI options
            const enableDPI = document.getElementById('enableDPI');
            const dpiValue = document.getElementById('dpiValue');
            const dpiOptions = document.getElementById('dpiOptions');
            
            if (enableDPI) {
                enableDPI.checked = options.enableDPI !== false; // Default true
            }
            if (dpiValue) {
                dpiValue.value = options.dpiValue || '300';
            }
            if (dpiOptions) {
                if (options.enableDPI !== false) {
                    dpiOptions.classList.remove('hidden');
                } else {
                    dpiOptions.classList.add('hidden');
                }
            }
        } catch (error) {
            console.error('Failed to load saved options:', error);
            setDefaultOptions();
        }
    } else {
        setDefaultOptions();
    }
}

// Set default options for first-time users
function setDefaultOptions() {
    const enableResize = document.getElementById('enableResize');
    const resizeOptions = document.getElementById('resizeOptions');
    const enableDPI = document.getElementById('enableDPI');
    const dpiOptions = document.getElementById('dpiOptions');
    const dpiValue = document.getElementById('dpiValue');
    
    if (enableResize) {
        enableResize.checked = true;
    }
    if (resizeOptions) {
        resizeOptions.classList.remove('hidden');
    }
    if (enableDPI) {
        enableDPI.checked = true;
    }
    if (dpiOptions) {
        dpiOptions.classList.remove('hidden');
    }
    if (dpiValue) {
        dpiValue.value = '300';
    }
    
    saveOptions();
}

// Save current options to localStorage
function saveOptions() {
    const enableResize = document.getElementById('enableResize');
    const resizeWidth = document.getElementById('resizeWidth');
    const resizeHeight = document.getElementById('resizeHeight');
    const enableDPI = document.getElementById('enableDPI');
    const dpiValue = document.getElementById('dpiValue');
    
    // Only save if elements exist
    if (!enableResize || !enableDPI) {
        return;
    }
    
    const options = {
        enableResize: enableResize.checked,
        resizeWidth: resizeWidth ? resizeWidth.value : '',
        resizeHeight: resizeHeight ? resizeHeight.value : '',
        enableDPI: enableDPI.checked,
        dpiValue: dpiValue ? dpiValue.value : '300'
    };
    
    localStorage.setItem('imagekitOptions', JSON.stringify(options));
}

// Handle file selection
function handleFileSelect(e) {
    handleFiles(e.target.files);
}

// Handle files
function handleFiles(files) {
    const validFiles = Array.from(files).filter(file => {
        if (!file.type.match(/^image\/(jpeg|png)$/)) {
            console.warn(`Skipping ${file.name}: not a JPEG or PNG`);
            return false;
        }
        if (file.size > 10 * 1024 * 1024) {
            console.warn(`Skipping ${file.name}: larger than 10MB`);
            return false;
        }
        return true;
    });
    
    if (validFiles.length === 0) {
        alert('유효한 이미지 파일이 없습니다.');
        return;
    }
    
    selectedFiles = validFiles;
    displaySelectedImages();
    
    document.getElementById('imageList').classList.remove('hidden');
    document.getElementById('optionsPanel').classList.remove('hidden');
}

// Display selected images
function displaySelectedImages() {
    const imageGrid = document.getElementById('imageGrid');
    imageGrid.innerHTML = '';
    
    selectedFiles.forEach((file, index) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            const item = document.createElement('div');
            item.className = 'relative bg-gray-100 rounded-lg overflow-hidden aspect-square';
            item.innerHTML = `
                <img src="${e.target.result}" alt="${file.name}" class="w-full h-full object-cover">
                <button class="absolute top-2 right-2 w-7 h-7 bg-red-600 hover:bg-red-700 text-white rounded-full flex items-center justify-center text-lg leading-none remove-btn" data-index="${index}">×</button>
            `;
            
            item.querySelector('.remove-btn').addEventListener('click', () => {
                selectedFiles.splice(index, 1);
                if (selectedFiles.length === 0) {
                    document.getElementById('imageList').classList.add('hidden');
                    document.getElementById('optionsPanel').classList.add('hidden');
                    document.getElementById('results').classList.add('hidden');
                } else {
                    displaySelectedImages();
                }
            });
            
            imageGrid.appendChild(item);
        };
        reader.readAsDataURL(file);
    });
}

// Process images
async function processImages() {
    if (!wasmReady) {
        alert('WebAssembly가 아직 로드되지 않았습니다.');
        return;
    }
    
    if (selectedFiles.length === 0) {
        alert('처리할 이미지를 선택해주세요.');
        return;
    }
    
    const options = getProcessingOptions();
    console.log('Processing options:', options);
    
    if (!options.resize && !options.crop && !options.dpi) {
        alert('최소 하나의 변환 옵션을 선택해주세요.');
        return;
    }
    
    document.getElementById('loading').classList.remove('hidden');
    // Change only the first p element (main message)
    const loadingMessages = document.querySelectorAll('#loading p');
    if (loadingMessages.length > 0) {
        loadingMessages[0].textContent = '이미지 처리 중...';
        // Hide the second message (about initial loading time)
        const initialLoadMessage = document.querySelector('.initial-load-message');
        if (initialLoadMessage) {
            initialLoadMessage.style.display = 'none';
        }
    }
    
    processedImages = [];
    
    for (const file of selectedFiles) {
        console.log(`Starting to process ${file.name}...`);
        try {
            const processedImage = await processImageFile(file, options);
            console.log(`Successfully processed ${file.name}:`, processedImage);
            
            if (processedImage && processedImage.data) {
                processedImages.push(processedImage);
                console.log(`Added ${file.name} to processed images`);
            } else {
                console.error(`Invalid result for ${file.name}:`, processedImage);
            }
        } catch (error) {
            console.error(`Failed to process ${file.name}:`, error);
            console.error('Error stack:', error.stack);
            alert(`이미지 처리 실패: ${file.name}\n${error.message}`);
        }
    }
    
    console.log('All processed images:', processedImages);
    document.getElementById('loading').classList.add('hidden');
    
    if (processedImages.length > 0) {
        displayResults();
    } else {
        alert('처리된 이미지가 없습니다. 콘솔을 확인해주세요.');
    }
}

// Get processing options
function getProcessingOptions() {
    const options = {
        resize: document.getElementById('enableResize').checked,
        crop: false, // Crop is now handled via individual image dialog
        dpi: document.getElementById('enableDPI').checked,
        width: 0,
        height: 0,
        cropTop: '0',
        cropRight: '0',
        cropBottom: '0',
        cropLeft: '0'
    };
    
    if (options.resize) {
        options.width = parseSize(document.getElementById('resizeWidth').value) || 0;
        options.height = parseSize(document.getElementById('resizeHeight').value) || 0;
    }
    
    if (options.dpi) {
        options.dpi = parseInt(document.getElementById('dpiValue').value) || 96;
    }
    
    return options;
}

// Parse size value (handles px and multipliers)
function parseSize(value) {
    if (!value) return 0;
    
    value = value.trim().toLowerCase();
    
    // Handle multipliers (2x, x2, 0.5x)
    if (value.includes('x')) {
        const multiplier = parseFloat(value.replace('x', ''));
        return isNaN(multiplier) ? 0 : -multiplier; // Negative for multiplier
    }
    
    // Handle pixel values
    const pixels = parseInt(value);
    return isNaN(pixels) ? 0 : pixels;
}

// Process single image
async function processImageFile(file, options) {
    console.log('processImageFile called with:', file.name, file.type, file.size, 'bytes');
    console.log('Options received:', options);
    
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        
        reader.onerror = (error) => {
            console.error('FileReader error:', error);
            reject(new Error('Failed to read file: ' + error));
        };
        
        reader.onabort = () => {
            console.error('FileReader aborted');
            reject(new Error('File reading was aborted'));
        };
        
        reader.onload = async (e) => {
            console.log('FileReader onload triggered');
            try {
                const imageDataUrl = e.target.result;
                console.log('Image loaded, data URL length:', imageDataUrl.length);
                console.log('Image data URL preview:', imageDataUrl.substring(0, 100));
                
                // Clone options to avoid modifying the original
                const processOptions = JSON.parse(JSON.stringify(options));
                
                // Get original dimensions if needed for multipliers
                if (processOptions.resize && (processOptions.width < 0 || processOptions.height < 0)) {
                    console.log('Multiplier detected, getting original dimensions...');
                    const img = new Image();
                    await new Promise((imgResolve, imgReject) => {
                        img.onload = () => {
                            console.log('Original dimensions:', img.width, 'x', img.height);
                            imgResolve();
                        };
                        img.onerror = () => {
                            imgReject(new Error('Failed to load image for dimension calculation'));
                        };
                        img.src = imageDataUrl;
                    });
                    
                    // Convert multipliers to actual pixels
                    if (processOptions.width < 0) {
                        processOptions.width = Math.round(img.width * Math.abs(processOptions.width));
                        console.log('Converted width multiplier to:', processOptions.width);
                    }
                    if (processOptions.height < 0) {
                        processOptions.height = Math.round(img.height * Math.abs(processOptions.height));
                        console.log('Converted height multiplier to:', processOptions.height);
                    }
                }
                
                // Call WebAssembly function
                if (typeof window.processImage !== 'function') {
                    throw new Error('WebAssembly function not loaded');
                }
                
                console.log('Calling WASM with data length:', imageDataUrl.length);
                console.log('Final options for WASM:', processOptions);
                
                const result = window.processImage(imageDataUrl, processOptions);
                console.log('WASM result:', result);
                
                if (!result || !result.success) {
                    throw new Error(result ? result.error : 'Unknown error');
                }
                
                if (!result.data) {
                    throw new Error('No image data returned');
                }
                
                resolve({
                    name: file.name,
                    originalSize: file.size,
                    data: result.data,
                    dpi: result.dpi || null
                });
            } catch (error) {
                reject(error);
            }
        };
        
        reader.readAsDataURL(file);
    });
}

// Display results
function displayResults() {
    const results = document.getElementById('results');
    const resultGrid = document.getElementById('resultGrid');
    
    resultGrid.innerHTML = '';
    
    processedImages.forEach((image, index) => {
        if (!image || !image.data) {
            console.error('Invalid image data at index', index);
            return;
        }
        
        const item = document.createElement('div');
        item.className = 'bg-gray-50 rounded-lg overflow-hidden';
        
        // Calculate new size from base64
        const base64Parts = image.data.split(',');
        const base64Length = base64Parts.length > 1 ? base64Parts[1].length : 0;
        const newSize = Math.round(base64Length * 0.75);
        
        // Create image element to get dimensions
        const img = item.querySelector('img') || document.createElement('img');
        
        item.innerHTML = `
            <img src="${image.data}" alt="${image.name}" class="w-full h-48 object-cover result-image">
            <div class="p-4">
                <p class="font-semibold text-gray-800 mb-2 truncate">${image.name}</p>
                <p class="text-sm text-gray-600">원본: ${formatFileSize(image.originalSize)}</p>
                <p class="text-sm text-gray-600">변환: ${formatFileSize(newSize)}</p>
                <p class="text-sm text-gray-600 dimensions-info">크기: 로딩중...</p>
                ${image.dpi ? `<p class="text-sm text-gray-600">DPI: ${image.dpi}</p>` : ''}
                <button class="w-full mt-3 mb-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors download-btn" data-index="${index}">다운로드</button>
                <button class="w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors crop-btn" data-index="${index}">크롭</button>
            </div>
        `;
        
        // Load image to get dimensions
        const resultImg = item.querySelector('.result-image');
        resultImg.onload = function() {
            const dimensionsInfo = item.querySelector('.dimensions-info');
            if (dimensionsInfo) {
                dimensionsInfo.textContent = `크기: ${this.naturalWidth} × ${this.naturalHeight}px`;
            }
        };
        
        item.querySelector('.download-btn').addEventListener('click', () => {
            downloadImage(image);
        });
        
        item.querySelector('.crop-btn').addEventListener('click', () => {
            openCropDialog(image);
        });
        
        resultGrid.appendChild(item);
    });
    
    results.classList.remove('hidden');
    results.scrollIntoView({ behavior: 'smooth' });
}

// Format file size
function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

// Download single image
function downloadImage(image) {
    if (!image || !image.data || !image.name) {
        console.error('Invalid image data for download:', image);
        return;
    }
    
    const link = document.createElement('a');
    link.href = image.data;
    
    // Generate filename
    const nameParts = image.name.split('.');
    const ext = nameParts.pop();
    const baseName = nameParts.join('.');
    link.download = `${baseName}_converted.${ext}`;
    
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

// Download all images
function downloadAll() {
    if (!processedImages || processedImages.length === 0) {
        alert('변환된 이미지가 없습니다.');
        return;
    }
    
    processedImages.forEach((image, index) => {
        if (image && image.data) {
            setTimeout(() => {
                downloadImage(image);
            }, index * 100); // Delay to prevent browser blocking
        }
    });
}

// Crop Dialog Class
class CropDialog {
    constructor() {
        this.modal = document.getElementById('cropModal');
        this.canvas = document.getElementById('cropCanvas');
        this.ctx = this.canvas.getContext('2d');
        this.previewCanvas = document.getElementById('cropPreview');
        this.previewCtx = this.previewCanvas.getContext('2d');
        this.coordsDiv = document.getElementById('cropCoords');
        
        this.image = null;
        this.imageData = null;
        this.imageName = null;
        this.isDrawing = false;
        this.startX = 0;
        this.startY = 0;
        this.endX = 0;
        this.endY = 0;
        this.ratio = 0; // 0 = free, >0 = fixed ratio
        
        this.setupEventListeners();
    }
    
    setupEventListeners() {
        // Modal controls
        document.getElementById('cropModalClose').addEventListener('click', () => this.close());
        document.getElementById('cropCancel').addEventListener('click', () => this.close());
        document.getElementById('cropApply').addEventListener('click', () => this.applyCrop());
        
        // Canvas events - mouse events
        this.canvas.addEventListener('mousedown', (e) => this.startDraw(e));
        
        // Canvas events - touch events for mobile/iOS
        this.canvas.addEventListener('touchstart', (e) => this.startDrawTouch(e), { passive: false });
        this.canvas.addEventListener('touchmove', (e) => this.drawTouch(e), { passive: false });
        this.canvas.addEventListener('touchend', (e) => this.endDrawTouch(e), { passive: false });
        this.canvas.addEventListener('touchcancel', (e) => this.endDrawTouch(e), { passive: false });
        
        // Document-level events for mousemove and mouseup to track outside canvas
        this.documentMouseMove = (e) => this.draw(e);
        this.documentMouseUp = () => this.endDraw();
        
        // Preset buttons
        document.querySelectorAll('.preset-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('.preset-btn').forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                this.ratio = parseFloat(btn.dataset.ratio);
                this.redrawSelection();
            });
        });
    }
    
    open(imageData, imageName) {
        this.imageData = imageData;
        this.imageName = imageName;
        this.modal.classList.remove('hidden');
        
        // Reset ratio and set first button as active
        this.ratio = 0;
        document.querySelectorAll('.preset-btn').forEach((btn, index) => {
            if (index === 0) {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });
        
        // Load image
        this.image = new Image();
        this.image.onload = () => {
            this.initCanvas();
        };
        this.image.src = imageData;
    }
    
    close() {
        this.modal.classList.add('hidden');
        this.clearSelection();
        
        // Ensure document listeners are removed
        document.removeEventListener('mousemove', this.documentMouseMove);
        document.removeEventListener('mouseup', this.documentMouseUp);
        this.isDrawing = false;
    }
    
    initCanvas() {
        // Calculate scale to fit canvas - use most of modal space
        const modalBody = document.getElementById('cropModalBody');
        const containerWidth = modalBody ? modalBody.clientWidth * 0.65 : window.innerWidth * 0.6;
        const containerHeight = window.innerHeight * 0.75;
        
        const maxWidth = Math.min(1200, containerWidth);
        const maxHeight = Math.min(800, containerHeight);
        
        let width = this.image.width;
        let height = this.image.height;
        
        if (width > maxWidth || height > maxHeight) {
            const scale = Math.min(maxWidth / width, maxHeight / height);
            width = Math.round(width * scale);
            height = Math.round(height * scale);
        }
        
        this.canvas.width = width;
        this.canvas.height = height;
        this.scale = width / this.image.width;
        
        // Draw image
        this.ctx.drawImage(this.image, 0, 0, width, height);
        
        // Clear any previous selection
        this.clearSelection();
    }
    
    startDraw(e) {
        e.preventDefault();
        const rect = this.canvas.getBoundingClientRect();
        this.isDrawing = true;
        this.canvasRect = rect; // Store for use in draw method
        
        // Get mouse position relative to canvas, accounting for any scroll
        this.startX = e.clientX - rect.left;
        this.startY = e.clientY - rect.top;
        
        // Clamp to canvas bounds
        this.startX = Math.max(0, Math.min(this.startX, this.canvas.width));
        this.startY = Math.max(0, Math.min(this.startY, this.canvas.height));
        
        this.endX = this.startX;
        this.endY = this.startY;
        
        // Add document-level listeners for tracking mouse outside canvas
        document.addEventListener('mousemove', this.documentMouseMove);
        document.addEventListener('mouseup', this.documentMouseUp);
    }
    
    draw(e) {
        if (!this.isDrawing) return;
        
        e.preventDefault();
        // Use stored canvas rect or get new one
        const rect = this.canvasRect || this.canvas.getBoundingClientRect();
        
        // Get mouse position relative to canvas
        this.endX = e.clientX - rect.left;
        this.endY = e.clientY - rect.top;
        
        // Clamp to canvas bounds - this allows tracking outside but constrains to edges
        this.endX = Math.max(0, Math.min(this.endX, this.canvas.width));
        this.endY = Math.max(0, Math.min(this.endY, this.canvas.height));
        
        // Apply ratio constraint if set
        if (this.ratio > 0) {
            const width = Math.abs(this.endX - this.startX);
            const height = width / this.ratio;
            this.endY = this.startY + (this.endY > this.startY ? height : -height);
            
            // Clamp again after ratio adjustment
            this.endY = Math.max(0, Math.min(this.endY, this.canvas.height));
        }
        
        this.redrawSelection();
        this.updatePreview();
        this.updateCoords();
    }
    
    endDraw() {
        this.isDrawing = false;
        
        // Remove document-level listeners
        document.removeEventListener('mousemove', this.documentMouseMove);
        document.removeEventListener('mouseup', this.documentMouseUp);
    }
    
    // Touch event handlers for mobile/iOS support
    startDrawTouch(e) {
        e.preventDefault();
        
        // Get the first touch point
        if (e.touches.length === 0) return;
        const touch = e.touches[0];
        
        const rect = this.canvas.getBoundingClientRect();
        this.isDrawing = true;
        this.canvasRect = rect;
        
        // Get touch position relative to canvas
        this.startX = touch.clientX - rect.left;
        this.startY = touch.clientY - rect.top;
        
        // Clamp to canvas bounds
        this.startX = Math.max(0, Math.min(this.startX, this.canvas.width));
        this.startY = Math.max(0, Math.min(this.startY, this.canvas.height));
        
        this.endX = this.startX;
        this.endY = this.startY;
    }
    
    drawTouch(e) {
        if (!this.isDrawing) return;
        
        e.preventDefault();
        
        // Get the first touch point
        if (e.touches.length === 0) return;
        const touch = e.touches[0];
        
        // Use stored canvas rect or get new one
        const rect = this.canvasRect || this.canvas.getBoundingClientRect();
        
        // Get touch position relative to canvas
        this.endX = touch.clientX - rect.left;
        this.endY = touch.clientY - rect.top;
        
        // Clamp to canvas bounds
        this.endX = Math.max(0, Math.min(this.endX, this.canvas.width));
        this.endY = Math.max(0, Math.min(this.endY, this.canvas.height));
        
        // Apply ratio constraint if set
        if (this.ratio > 0) {
            const width = Math.abs(this.endX - this.startX);
            const height = width / this.ratio;
            this.endY = this.startY + (this.endY > this.startY ? height : -height);
            
            // Clamp again after ratio adjustment
            this.endY = Math.max(0, Math.min(this.endY, this.canvas.height));
        }
        
        this.redrawSelection();
        this.updatePreview();
        this.updateCoords();
    }
    
    endDrawTouch(e) {
        e.preventDefault();
        this.isDrawing = false;
    }
    
    redrawSelection() {
        // Clear and redraw image
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.drawImage(this.image, 0, 0, this.canvas.width, this.canvas.height);
        
        if (Math.abs(this.endX - this.startX) < 5 || Math.abs(this.endY - this.startY) < 5) {
            return;
        }
        
        // Draw selection rectangle
        const x = Math.min(this.startX, this.endX);
        const y = Math.min(this.startY, this.endY);
        const width = Math.abs(this.endX - this.startX);
        const height = Math.abs(this.endY - this.startY);
        
        // Semi-transparent overlay
        this.ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Clear selection area
        this.ctx.clearRect(x, y, width, height);
        this.ctx.drawImage(this.image, 
            x / this.scale, y / this.scale, width / this.scale, height / this.scale,
            x, y, width, height);
        
        // Draw border
        this.ctx.strokeStyle = '#667eea';
        this.ctx.lineWidth = 2;
        this.ctx.setLineDash([5, 5]);
        this.ctx.strokeRect(x, y, width, height);
        this.ctx.setLineDash([]);
    }
    
    updatePreview() {
        const x = Math.min(this.startX, this.endX);
        const y = Math.min(this.startY, this.endY);
        const width = Math.abs(this.endX - this.startX);
        const height = Math.abs(this.endY - this.startY);
        
        if (width < 5 || height < 5) return;
        
        // Set preview canvas size
        const maxPreviewSize = 180;
        let previewWidth = width;
        let previewHeight = height;
        
        if (previewWidth > maxPreviewSize || previewHeight > maxPreviewSize) {
            const scale = Math.min(maxPreviewSize / previewWidth, maxPreviewSize / previewHeight);
            previewWidth = Math.round(previewWidth * scale);
            previewHeight = Math.round(previewHeight * scale);
        }
        
        this.previewCanvas.width = previewWidth;
        this.previewCanvas.height = previewHeight;
        
        // Draw preview
        this.previewCtx.drawImage(this.image,
            x / this.scale, y / this.scale, width / this.scale, height / this.scale,
            0, 0, previewWidth, previewHeight);
    }
    
    updateCoords() {
        const x = Math.min(this.startX, this.endX);
        const y = Math.min(this.startY, this.endY);
        const width = Math.abs(this.endX - this.startX);
        const height = Math.abs(this.endY - this.startY);
        
        const realX = Math.round(x / this.scale);
        const realY = Math.round(y / this.scale);
        const realWidth = Math.round(width / this.scale);
        const realHeight = Math.round(height / this.scale);
        
        this.coordsDiv.innerHTML = `
            위치: ${realX}, ${realY}<br>
            크기: ${realWidth} × ${realHeight}px
        `;
    }
    
    clearSelection() {
        this.startX = 0;
        this.startY = 0;
        this.endX = 0;
        this.endY = 0;
        this.coordsDiv.innerHTML = '';
        if (this.image) {
            this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
            this.ctx.drawImage(this.image, 0, 0, this.canvas.width, this.canvas.height);
        }
        this.previewCtx.clearRect(0, 0, this.previewCanvas.width, this.previewCanvas.height);
    }
    
    async applyCrop() {
        const x = Math.min(this.startX, this.endX);
        const y = Math.min(this.startY, this.endY);
        const width = Math.abs(this.endX - this.startX);
        const height = Math.abs(this.endY - this.startY);
        
        if (width < 5 || height < 5) {
            alert('크롭 영역을 선택해주세요.');
            return;
        }
        
        // Calculate crop parameters for WASM
        const realX = Math.round(x / this.scale);
        const realY = Math.round(y / this.scale);
        const realWidth = Math.round(width / this.scale);
        const realHeight = Math.round(height / this.scale);
        
        try {
            // Call WebAssembly crop function
            const cropOptions = {
                resize: false,
                crop: true,
                dpi: false,
                width: 0,
                height: 0,
                cropTop: realY.toString(),
                cropRight: (this.image.width - realX - realWidth).toString(),
                cropBottom: (this.image.height - realY - realHeight).toString(),
                cropLeft: realX.toString()
            };
            
            console.log('Crop options:', cropOptions);
            
            const result = window.processImage(this.imageData, cropOptions);
            
            if (result.success && result.data) {
                // Download cropped image
                const croppedImage = {
                    name: this.imageName.replace(/\.(jpg|jpeg|png)$/i, '_cropped.$1'),
                    data: result.data
                };
                downloadImage(croppedImage);
                this.close();
            } else {
                alert('크롭 처리 실패: ' + (result.error || 'Unknown error'));
            }
        } catch (error) {
            console.error('Crop error:', error);
            alert('크롭 처리 중 오류가 발생했습니다.');
        }
    }
}

// Initialize crop dialog
let cropDialog;

// Open crop dialog
function openCropDialog(image) {
    if (!cropDialog) {
        cropDialog = new CropDialog();
    }
    cropDialog.open(image.data, image.name);
}