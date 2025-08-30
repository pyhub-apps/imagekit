// Global variables
let wasmReady = false;
let selectedFiles = [];
let processedImages = [];

// Initialize WebAssembly
async function initWasm() {
    const go = new Go();
    
    try {
        console.log('Loading WebAssembly...');
        const response = await fetch('static/imagekit.wasm');
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

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initWasm();
    setupEventListeners();
});

// Setup event listeners
function setupEventListeners() {
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    
    // Load saved options from localStorage
    loadSavedOptions();
    
    // File upload
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
    
    // Options toggles with auto-save
    document.getElementById('enableResize').addEventListener('change', (e) => {
        document.getElementById('resizeOptions').style.display = e.target.checked ? 'block' : 'none';
        saveOptions();
    });
    
    document.getElementById('enableCrop').addEventListener('change', (e) => {
        document.getElementById('cropOptions').style.display = e.target.checked ? 'block' : 'none';
        saveOptions();
    });
    
    document.getElementById('enableDPI').addEventListener('change', (e) => {
        document.getElementById('dpiOptions').style.display = e.target.checked ? 'block' : 'none';
        saveOptions();
    });
    
    // Save options on input change
    document.getElementById('resizeWidth').addEventListener('input', saveOptions);
    document.getElementById('resizeHeight').addEventListener('input', saveOptions);
    document.getElementById('cropTop').addEventListener('input', saveOptions);
    document.getElementById('cropRight').addEventListener('input', saveOptions);
    document.getElementById('cropBottom').addEventListener('input', saveOptions);
    document.getElementById('cropLeft').addEventListener('input', saveOptions);
    document.getElementById('dpiValue').addEventListener('change', saveOptions);
    
    // Process button
    document.getElementById('processBtn').addEventListener('click', processImages);
    
    // Download all button
    document.getElementById('downloadAllBtn').addEventListener('click', downloadAll);
}

// Load saved options from localStorage
function loadSavedOptions() {
    const saved = localStorage.getItem('imagekitOptions');
    
    if (saved) {
        try {
            const options = JSON.parse(saved);
            
            // Load resize options
            document.getElementById('enableResize').checked = options.enableResize !== false; // Default true
            document.getElementById('resizeWidth').value = options.resizeWidth || '';
            document.getElementById('resizeHeight').value = options.resizeHeight || '';
            document.getElementById('resizeOptions').style.display = options.enableResize !== false ? 'block' : 'none';
            
            // Load crop options
            document.getElementById('enableCrop').checked = options.enableCrop || false;
            document.getElementById('cropTop').value = options.cropTop || '';
            document.getElementById('cropRight').value = options.cropRight || '';
            document.getElementById('cropBottom').value = options.cropBottom || '';
            document.getElementById('cropLeft').value = options.cropLeft || '';
            document.getElementById('cropOptions').style.display = options.enableCrop ? 'block' : 'none';
            
            // Load DPI options
            document.getElementById('enableDPI').checked = options.enableDPI !== false; // Default true
            document.getElementById('dpiValue').value = options.dpiValue || '300';
            document.getElementById('dpiOptions').style.display = options.enableDPI !== false ? 'block' : 'none';
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
    document.getElementById('enableResize').checked = true;
    document.getElementById('resizeOptions').style.display = 'block';
    document.getElementById('enableDPI').checked = true;
    document.getElementById('dpiOptions').style.display = 'block';
    document.getElementById('dpiValue').value = '300';
    saveOptions();
}

// Save current options to localStorage
function saveOptions() {
    const options = {
        enableResize: document.getElementById('enableResize').checked,
        resizeWidth: document.getElementById('resizeWidth').value,
        resizeHeight: document.getElementById('resizeHeight').value,
        enableCrop: document.getElementById('enableCrop').checked,
        cropTop: document.getElementById('cropTop').value,
        cropRight: document.getElementById('cropRight').value,
        cropBottom: document.getElementById('cropBottom').value,
        cropLeft: document.getElementById('cropLeft').value,
        enableDPI: document.getElementById('enableDPI').checked,
        dpiValue: document.getElementById('dpiValue').value
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
    
    document.getElementById('imageList').style.display = 'block';
    document.getElementById('optionsPanel').style.display = 'block';
}

// Display selected images
function displaySelectedImages() {
    const imageGrid = document.getElementById('imageGrid');
    imageGrid.innerHTML = '';
    
    selectedFiles.forEach((file, index) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            const item = document.createElement('div');
            item.className = 'image-item';
            item.innerHTML = `
                <img src="${e.target.result}" alt="${file.name}">
                <button class="remove-btn" data-index="${index}">×</button>
            `;
            
            item.querySelector('.remove-btn').addEventListener('click', () => {
                selectedFiles.splice(index, 1);
                if (selectedFiles.length === 0) {
                    document.getElementById('imageList').style.display = 'none';
                    document.getElementById('optionsPanel').style.display = 'none';
                    document.getElementById('results').style.display = 'none';
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
    document.querySelector('.loading p').textContent = '이미지 처리 중...';
    
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
        crop: document.getElementById('enableCrop').checked,
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
    
    if (options.crop) {
        options.cropTop = document.getElementById('cropTop').value || '0';
        options.cropRight = document.getElementById('cropRight').value || '0';
        options.cropBottom = document.getElementById('cropBottom').value || '0';
        options.cropLeft = document.getElementById('cropLeft').value || '0';
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
        item.className = 'result-item';
        
        // Calculate new size from base64
        const base64Parts = image.data.split(',');
        const base64Length = base64Parts.length > 1 ? base64Parts[1].length : 0;
        const newSize = Math.round(base64Length * 0.75);
        
        item.innerHTML = `
            <img src="${image.data}" alt="${image.name}">
            <div class="result-info">
                <p><strong>${image.name}</strong></p>
                <p>원본: ${formatFileSize(image.originalSize)}</p>
                <p>변환: ${formatFileSize(newSize)}</p>
                ${image.dpi ? `<p>DPI: ${image.dpi}</p>` : ''}
                <button class="download-btn" data-index="${index}">다운로드</button>
            </div>
        `;
        
        item.querySelector('.download-btn').addEventListener('click', () => {
            downloadImage(image);
        });
        
        resultGrid.appendChild(item);
    });
    
    results.style.display = 'block';
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