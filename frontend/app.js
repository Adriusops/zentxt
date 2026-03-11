import {
    ListFiles,
    CreateFile,
    SaveVersion,
    ListVersions,
    GetFile,
    GetVersion,
    RestoreVersion,
    DeleteVersion,
    ExportVersion
} from "./bindings/github.com/Adriusops/zentxt/cmd/zentxt/app.js";
import * as Dialog from "./node_modules/@wailsio/runtime/dist/dialogs.js";

// ======================
// ROUTER
// ======================

let currentFileId = null;
let currentFile = null;
let currentVersions = [];

function initRouter() {
    // Handle hash changes
    window.addEventListener('hashchange', handleRoute);

    // Handle initial route
    handleRoute();
}

function handleRoute() {
    const hash = window.location.hash.slice(1); // Remove #
    const [route, queryString] = hash.split('?');

    if (route === 'timeline' && queryString) {
        const params = new URLSearchParams(queryString);
        const fileId = params.get('fileId');
        if (fileId) {
            showTimeline(fileId);
        } else {
            showHome();
        }
    } else {
        showHome();
    }
}

function showHome() {
    document.getElementById('home-view').classList.remove('hidden');
    document.getElementById('timeline-view').classList.add('hidden');
    loadFiles();
}

function showTimeline(fileId) {
    currentFileId = fileId;
    document.getElementById('home-view').classList.add('hidden');
    document.getElementById('timeline-view').classList.remove('hidden');
    loadTimeline(fileId);
}

function navigateTo(route) {
    window.location.hash = route;
}

// ======================
// HOME VIEW
// ======================

// Format date to relative format (Today, Yesterday, or full date)
function formatRelativeDate(isoString) {
    const date = new Date(isoString);
    const now = new Date();
    const diffTime = now - date;
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

    const timeStr = date.toLocaleTimeString("fr-FR", {
        hour: "2-digit",
        minute: "2-digit",
    });

    if (diffDays === 0) {
        return `Today, ${timeStr}`;
    } else if (diffDays === 1) {
        return `Yesterday, ${timeStr}`;
    } else {
        return (
            date.toLocaleDateString("fr-FR", {
                day: "2-digit",
                month: "2-digit",
                year: "numeric",
            }) +
            ", " +
            timeStr
        );
    }
}

// Load and display files
async function loadFiles() {
    try {
        const files = await ListFiles();
        const filesList = document.getElementById("files-list");

        if (!files || files.length === 0) {
            filesList.innerHTML = `
                <div class="col-span-full text-center py-12">
                    <p class="text-gray-400 text-sm">
                        No files yet. Drop one above to get started.
                    </p>
                </div>
            `;
            return;
        }

        filesList.innerHTML = "";

        for (const file of files) {
            if (!file) continue;

            const card = createFileCard(file);
            filesList.appendChild(card);

            // Load preview for this file
            loadFilePreview(file.ID, card);
        }

        // Setup drop zones on cards
        setupFileCardDropZones();
    } catch (error) {
        console.error("Error loading files:", error);
        const filesList = document.getElementById("files-list");
        filesList.innerHTML = `
            <div class="col-span-full text-center py-12">
                <p class="text-red-500 text-sm">
                    Failed to load files. Please try again.
                </p>
            </div>
        `;
    }
}

// Create file card HTML element
function createFileCard(file) {
    const card = document.createElement("a");
    card.className = "file-card block group";
    card.dataset.fileId = file.ID;
    card.href = `#timeline?fileId=${file.ID}`;

    card.innerHTML = `
        <div class="file-drop-zone relative bg-white rounded-2xl overflow-hidden hover:shadow-xl transition-all border border-gray-200">
            <div class="bg-gray-50 p-8 pb-0 flex flex-col justify-end min-h-[240px] relative">
                <div class="mx-auto bg-white rounded-t-lg shadow-md overflow-hidden min-h-[200px] p-8 pb-12 relative border border-b-0 border-gray-300" style="max-width: 98%">
                    <div class="absolute inset-0 pointer-events-none opacity-20">
                        <div class="h-full flex flex-col justify-start pt-8 gap-[22px]">
                            <div class="w-full h-[1px] bg-gray-400"></div>
                            <div class="w-full h-[1px] bg-gray-400"></div>
                            <div class="w-full h-[1px] bg-gray-400"></div>
                            <div class="w-full h-[1px] bg-gray-400"></div>
                            <div class="w-full h-[1px] bg-gray-400"></div>
                            <div class="w-full h-[1px] bg-gray-400"></div>
                        </div>
                    </div>
                    <div class="preview-content text-sm text-gray-900 line-clamp-4 leading-relaxed font-body relative z-10" style="line-height: 24px">
                        <span class="text-gray-400 italic">Loading preview...</span>
                    </div>
                </div>
            </div>
            <div class="bg-white px-8 py-6 relative z-20 -mt-8">
                <h3 class="text-base font-bold text-gray-900 mb-2 font-display">
                    ${file.Name}
                </h3>
                <p class="text-sm text-gray-400 timestamp">
                    ${formatRelativeDate(file.CreatedAt)}
                </p>
            </div>
            <div class="drop-overlay absolute inset-0 bg-white bg-opacity-80 backdrop-blur-sm flex items-center justify-center opacity-0 pointer-events-none transition-opacity z-30">
                <div class="text-center">
                    <svg class="mx-auto h-12 w-12 text-gray-900 mb-3" stroke="currentColor" fill="none" viewBox="0 0 48 48">
                        <path d="M24 8v24m0 0l-8-8m8 8l8-8" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                    <p class="text-gray-900 font-medium text-sm">
                        Drop to save new version
                    </p>
                </div>
            </div>
        </div>
    `;

    return card;
}

// Load preview for a file
async function loadFilePreview(fileId, card) {
    const previewElement = card.querySelector(".preview-content");

    try {
        const versions = await ListVersions(fileId);

        if (versions && versions.length > 0) {
            const latestVersion = versions[0];
            const content = latestVersion?.Content || "";

            if (content) {
                const lines = content.split("\n").slice(0, 4);
                const preview = lines.join("\n").substring(0, 200);
                previewElement.textContent = preview + (content.length > 200 ? "..." : "");
                previewElement.classList.remove("text-gray-400", "italic");
            } else {
                previewElement.innerHTML = '<span class="text-gray-400 italic">Empty file</span>';
            }
        } else {
            previewElement.innerHTML = '<span class="text-gray-400 italic">No versions yet</span>';
        }
    } catch (error) {
        console.error("Error loading preview:", error);
        previewElement.innerHTML = '<span class="text-gray-400 italic">Failed to load preview</span>';
    }
}

// Drag & Drop functionality
function initHomeDragDrop() {
    const dropZone = document.getElementById("drop-zone");
    const fileInput = document.getElementById("file-input");

    // Click to browse
    dropZone.addEventListener("click", () => {
        fileInput.click();
    });

    // Prevent default drag behaviors
    ["dragenter", "dragover", "dragleave", "drop"].forEach((eventName) => {
        dropZone.addEventListener(eventName, preventDefaults, false);
        document.body.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    // Highlight drop zone when dragging
    ["dragenter", "dragover"].forEach((eventName) => {
        dropZone.addEventListener(eventName, () => {
            dropZone.classList.add("border-gray-400", "bg-gray-50");
        });
    });

    ["dragleave", "drop"].forEach((eventName) => {
        dropZone.addEventListener(eventName, () => {
            dropZone.classList.remove("border-gray-400", "bg-gray-50");
        });
    });

    // Handle dropped files
    dropZone.addEventListener("drop", (e) => {
        const files = e.dataTransfer.files;
        handleFiles(files);
    });

    // Handle file input change
    fileInput.addEventListener("change", (e) => {
        const files = e.target.files;
        handleFiles(files);
    });
}

// Modal management
function initHomeModal() {
    const modal = document.getElementById("version-modal");
    const versionForm = document.getElementById("version-form");
    const versionMessageInput = document.getElementById("version-message");
    const cancelButton = document.getElementById("cancel-version");

    let pendingVersionData = null;

    function showModal() {
        modal.classList.remove("hidden");
        versionMessageInput.value = "";
        versionMessageInput.focus();
    }

    function hideModal() {
        modal.classList.add("hidden");
        pendingVersionData = null;
    }

    cancelButton.addEventListener("click", hideModal);

    modal.addEventListener("click", (e) => {
        if (e.target === modal) {
            hideModal();
        }
    });

    // Handle form submission
    versionForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        if (!pendingVersionData) return;

        const message = versionMessageInput.value.trim();
        if (!message) return;

        const { fileId, content, fileName } = pendingVersionData;

        try {
            await SaveVersion(fileId, fileName, "You", message, content);
            navigateTo(`timeline?fileId=${fileId}`);
        } catch (error) {
            console.error("Error:", error);
            alert("Failed to save new version. Please try again.");
            hideModal();
        }
    });

    // Expose to global scope for file card drop zones
    window.showVersionModal = function(data) {
        pendingVersionData = data;
        showModal();
    };
}

// Drag & Drop on file cards to create new version
function setupFileCardDropZones() {
    const fileCards = document.querySelectorAll(".file-card");

    fileCards.forEach((card) => {
        const dropZoneEl = card.querySelector(".file-drop-zone");
        const overlay = card.querySelector(".drop-overlay");

        ["dragenter", "dragover", "dragleave", "drop"].forEach((eventName) => {
            dropZoneEl.addEventListener(
                eventName,
                (e) => {
                    e.preventDefault();
                    e.stopPropagation();
                },
                false,
            );
        });

        let dragCounter = 0;

        dropZoneEl.addEventListener("dragenter", (e) => {
            dragCounter++;
            overlay.classList.remove("opacity-0");
            overlay.classList.add("opacity-100", "pointer-events-auto");
        });

        dropZoneEl.addEventListener("dragover", (e) => {
            overlay.classList.remove("opacity-0");
            overlay.classList.add("opacity-100", "pointer-events-auto");
        });

        dropZoneEl.addEventListener("dragleave", (e) => {
            dragCounter--;
            if (dragCounter === 0) {
                overlay.classList.add("opacity-0");
                overlay.classList.remove("opacity-100", "pointer-events-auto");
            }
        });

        dropZoneEl.addEventListener("drop", (e) => {
            dragCounter = 0;
            overlay.classList.add("opacity-0");
            overlay.classList.remove("opacity-100", "pointer-events-auto");

            const files = e.dataTransfer.files;
            if (files.length === 0) return;

            const file = files[0];
            const fileId = card.dataset.fileId;

            const fileName = file.name.toLowerCase();
            if (!fileName.endsWith(".txt") && !fileName.endsWith(".md")) {
                alert(
                    "ZenTxt only supports .txt and .md files.\\n\\nPlease convert your document to plain text before uploading.",
                );
                return;
            }

            const reader = new FileReader();
            reader.onload = function (e) {
                const content = e.target.result;

                window.showVersionModal({
                    fileId: fileId,
                    content: content,
                    fileName: file.name,
                });
            };
            reader.readAsText(file);
        });
    });
}

// Process uploaded file
async function handleFiles(files) {
    if (files.length === 0) return;

    const file = files[0];

    const fileName = file.name.toLowerCase();
    if (!fileName.endsWith(".txt") && !fileName.endsWith(".md")) {
        alert(
            "ZenTxt only supports .txt and .md files.\\n\\nPlease convert your document to plain text before uploading.",
        );
        return;
    }

    const reader = new FileReader();

    reader.onload = async function (e) {
        const content = e.target.result;

        try {
            const newFile = await CreateFile(file.name, file.name);

            if (newFile) {
                await SaveVersion(newFile.ID, file.name, "You", "Initial version", content);
                loadFiles(); // Reload the file list
            }
        } catch (error) {
            console.error("Error:", error);
            alert("Failed to add file. Please try again.");
        }
    };

    reader.readAsText(file);
}

// ======================
// TIMELINE VIEW
// ======================

async function loadTimeline(fileId) {
    try {
        const file = await GetFile(fileId);
        const versions = await ListVersions(fileId);

        if (!file) {
            navigateTo('home');
            return;
        }

        currentFile = file;
        currentVersions = versions || [];

        document.getElementById('file-name').textContent = file.Name;

        renderTimeline();
    } catch (error) {
        console.error('Error loading timeline:', error);
        await Dialog.Error({
            title: 'Error',
            message: 'Failed to load file timeline. Please try again.'
        });
        navigateTo('home');
    }
}

function renderTimeline() {
    const container = document.getElementById('timeline-container');

    if (!currentVersions || currentVersions.length === 0) {
        container.innerHTML = `
            <div class="text-center py-12 border border-dashed border-gray-200 rounded-xl">
                <p class="text-gray-400 text-sm">No versions yet. Save your first version above.</p>
            </div>
        `;
        return;
    }

    container.innerHTML = currentVersions.map((version, index) => {
        if (!version) return '';

        const isLatest = index === 0;
        const isLast = index === currentVersions.length - 1;
        const prevVersionID = index < currentVersions.length - 1 ? currentVersions[index + 1]?.ID : null;

        return `
            <div class="timeline-item flex gap-4 relative" data-version-id="${version.ID}">
                <div class="timeline-dot-column flex flex-col items-center relative" style="width: 12px;">
                    <div class="timeline-dot-wrapper flex-shrink-0 relative z-10">
                        <div class="w-3 h-3 rounded-full ${isLatest ? 'bg-green-500' : 'bg-gray-300'}"></div>
                    </div>
                    ${!isLast ? '<div class="timeline-pipe"></div>' : ''}
                </div>

                <div class="flex-1 pb-8">
                    <div class="border border-gray-100 rounded-xl p-5 bg-white hover:border-gray-300 hover:shadow-md transition-all">
                        <div class="flex items-start justify-between mb-3">
                            <div class="flex-1">
                                <div class="flex items-center gap-2 mb-2">
                                    <span class="inline-block px-2.5 py-1 bg-gray-100 text-gray-700 text-xs font-semibold rounded">
                                        Version ${version.VersionNumber}
                                    </span>
                                    ${isLatest ? '<span class="inline-block px-2.5 py-1 bg-green-100 text-green-700 text-xs font-semibold rounded">Latest</span>' : ''}
                                </div>
                                <p class="text-base font-medium text-[#232323] mb-1">${version.Message}</p>
                                <p class="text-xs text-gray-500">${formatRelativeDate(version.CreatedAt)}</p>
                            </div>

                            <div class="relative version-menu-container">
                                <button
                                    data-action="toggle-menu"
                                    data-version-id="${version.ID}"
                                    class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                                    aria-label="More options"
                                >
                                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                                        <circle cx="12" cy="5" r="2"/>
                                        <circle cx="12" cy="12" r="2"/>
                                        <circle cx="12" cy="19" r="2"/>
                                    </svg>
                                </button>

                                <div
                                    id="menu-${version.ID}"
                                    class="version-dropdown hidden absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-xl border border-gray-200 z-50"
                                >
                                    <div class="py-1">
                                        <button
                                            data-action="share"
                                            data-version-id="${version.ID}"
                                            data-version-number="${version.VersionNumber}"
                                            data-version-message="${version.Message.replace(/"/g, '&quot;')}"
                                            class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-3 transition-colors"
                                        >
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                                            </svg>
                                            Share by email
                                        </button>

                                        <button
                                            data-action="download"
                                            data-version-id="${version.ID}"
                                            data-version-number="${version.VersionNumber}"
                                            class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-3 transition-colors"
                                        >
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
                                            </svg>
                                            Export version
                                        </button>

                                        <div class="border-t border-gray-100 my-1"></div>

                                        <button
                                            data-action="view"
                                            data-version-id="${version.ID}"
                                            class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-3 transition-colors"
                                        >
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                                            </svg>
                                            View
                                        </button>

                                        <div class="border-t border-gray-100 my-1"></div>

                                        <button
                                            data-action="delete"
                                            data-version-id="${version.ID}"
                                            data-version-number="${version.VersionNumber}"
                                            class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-3 transition-colors"
                                        >
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                                            </svg>
                                            Delete
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="flex flex-wrap gap-2 pt-3 border-t border-gray-100">
                            ${prevVersionID ? `<button
                                data-action="compare"
                                data-v1="${prevVersionID}"
                                data-v2="${version.ID}"
                                class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">
                                Compare
                            </button>` : ''}
                            ${!isLatest ? `<button
                                data-action="restore"
                                data-version-id="${version.ID}"
                                class="px-3 py-1.5 text-sm font-medium text-[#232323] bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                                Restore
                            </button>` : ''}
                        </div>
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

// Timeline event handlers
function initTimelineEventHandlers() {
    // Event delegation for all button clicks
    document.addEventListener('click', (e) => {
        const button = e.target.closest('button[data-action]');

        if (!button) {
            // Close all menus when clicking outside
            const allMenus = document.querySelectorAll('.version-dropdown');
            allMenus.forEach(m => m.classList.add('hidden'));
            return;
        }

        e.preventDefault();
        e.stopPropagation();

        const action = button.dataset.action;

        // Handle different actions
        switch(action) {
            case 'toggle-menu':
                const versionId = button.dataset.versionId;
                const menu = document.getElementById(`menu-${versionId}`);
                const allMenus = document.querySelectorAll('.version-dropdown');
                allMenus.forEach(m => {
                    if (m.id !== `menu-${versionId}`) {
                        m.classList.add('hidden');
                    }
                });
                menu.classList.toggle('hidden');
                break;

            case 'share':
                openShareModal(
                    button.dataset.versionId,
                    button.dataset.versionNumber,
                    button.dataset.versionMessage
                );
                break;

            case 'download':
                downloadVersion(
                    button.dataset.versionId,
                    parseInt(button.dataset.versionNumber)
                );
                break;

            case 'view':
                viewVersion(button.dataset.versionId);
                break;

            case 'delete':
                confirmDeleteVersion(
                    button.dataset.versionId,
                    parseInt(button.dataset.versionNumber)
                );
                break;

            case 'compare':
                compareVersions(button.dataset.v1, button.dataset.v2);
                break;

            case 'restore':
                restoreVersionClick(button.dataset.versionId);
                break;

            default:
                console.warn('Unknown action:', action);
        }
    });
}

// View version
async function viewVersion(versionId) {
    try {
        const version = await GetVersion(versionId);
        const content = version?.Content || '';

        const versionContent = document.getElementById('version-content');
        if (!content) {
            versionContent.textContent = 'No content available for this version.';
        } else {
            versionContent.textContent = content;
        }

        document.getElementById('version-modal').classList.remove('hidden');
    } catch (error) {
        console.error('Error:', error);
        await Dialog.Error({
            title: 'Error',
            message: 'Failed to load version. Please try again.'
        });
    }
}

// Download version
async function downloadVersion(versionId, versionNumber) {
    try {
        const extension = currentFile.Name.match(/\.[^/.]+$/);
        const ext = extension ? extension[0] : '.txt';
        const nameWithoutExt = currentFile.Name.replace(/\.[^/.]+$/, '');
        const defaultFilename = `${nameWithoutExt}-v${versionNumber}${ext}`;

        const filePath = await Dialog.SaveFile({
            title: 'Export version',
            filename: defaultFilename,
            defaultFilename: defaultFilename,
            defaultFileName: defaultFilename,
            filters: [
                { displayName: 'Text files', pattern: '*.txt' },
                { displayName: 'Markdown files', pattern: '*.md' },
                { displayName: 'All files', pattern: '*.*' }
            ]
        });

        if (!filePath) return;

        await ExportVersion(versionId, filePath);

        await Dialog.Info({
            title: 'Export successful',
            message: `Version exported successfully to:\n${filePath}`
        });
    } catch (error) {
        console.error('Error downloading version:', error);
        await Dialog.Error({
            title: 'Export failed',
            message: 'Failed to export version. Please try again.\n\nError: ' + error.message
        });
    }
}

// Compare versions
function compareVersions(v1, v2) {
    window.location.href = `diff.html?fileId=${currentFileId}&v1=${v1}&v2=${v2}`;
}

// Restore version
async function restoreVersionClick(versionId) {
    try {
        const response = await Dialog.Question({
            title: 'Restore version',
            message: 'Are you sure you want to restore this version? This will create a new version with this content.',
            buttons: [
                { label: 'Cancel', isDefault: true },
                { label: 'Restore' }
            ]
        });

        if (response !== 'Restore') return;

        const result = await RestoreVersion(currentFileId, versionId);

        if (result) {
            await Dialog.Info({
                title: 'Version restored',
                message: 'The version has been restored successfully.'
            });
            loadTimeline(currentFileId);
        } else {
            await Dialog.Error({
                title: 'Restore failed',
                message: 'Failed to restore version. The operation returned no result.'
            });
        }
    } catch (error) {
        console.error('Error restoring version:', error);
        await Dialog.Error({
            title: 'Restore failed',
            message: 'Failed to restore version. Please try again.\n\nError: ' + error.message
        });
    }
}

// Delete version
let pendingDelete = null;

function confirmDeleteVersion(versionId, versionNumber) {
    pendingDelete = { versionId, versionNumber };

    const totalVersions = currentVersions.length;
    const deleteWarning = document.getElementById('delete-warning');

    if (totalVersions === 1) {
        deleteWarning.classList.remove('hidden');
    } else {
        deleteWarning.classList.add('hidden');
    }

    document.getElementById('delete-modal').classList.remove('hidden');
}

function initTimelineModals() {
    // Delete modal
    document.getElementById('cancel-delete').addEventListener('click', () => {
        document.getElementById('delete-modal').classList.add('hidden');
        pendingDelete = null;
    });

    document.getElementById('delete-modal').addEventListener('click', (e) => {
        if (e.target.id === 'delete-modal') {
            document.getElementById('delete-modal').classList.add('hidden');
            pendingDelete = null;
        }
    });

    document.getElementById('confirm-delete').addEventListener('click', async () => {
        if (!pendingDelete) return;

        const { versionId } = pendingDelete;
        const totalVersions = currentVersions.length;
        const isLastVersion = totalVersions === 1;

        const timelineItem = document.querySelector(`[data-version-id="${versionId}"]`);

        document.getElementById('delete-modal').classList.add('hidden');

        if (isLastVersion) {
            try {
                await DeleteVersion(versionId);
                navigateTo('home');
            } catch (error) {
                console.error('Error:', error);
                await Dialog.Error({
                    title: 'Error',
                    message: 'Failed to delete version. Please try again.'
                });
            }
        } else {
            if (timelineItem) {
                timelineItem.classList.add('deleting');
            }

            const deletePromise = DeleteVersion(versionId);

            setTimeout(async () => {
                try {
                    await deletePromise;
                    loadTimeline(currentFileId);
                } catch (error) {
                    console.error('Error:', error);
                    await Dialog.Error({
                        title: 'Error',
                        message: 'Failed to delete version. Please try again.'
                    });
                    if (timelineItem) {
                        timelineItem.classList.remove('deleting');
                    }
                }
            }, 400);
        }
    });

    // Share modal
    let pendingShare = null;

    window.openShareModal = async function(versionId, versionNumber, versionMessage) {
        try {
            const version = await GetVersion(versionId);
            const content = version?.Content || '';

            pendingShare = {
                versionNumber,
                versionMessage,
                content
            };

            document.getElementById('share-version-number').textContent = versionNumber;
            document.getElementById('share-file-name').textContent = currentFile.Name;
            document.getElementById('share-version-info').textContent = versionMessage;

            document.getElementById('share-modal').classList.remove('hidden');
            document.getElementById('share-email').focus();
        } catch (error) {
            console.error('Error:', error);
            await Dialog.Error({
                title: 'Error',
                message: 'Failed to load version content. Please try again.'
            });
        }
    };

    document.getElementById('cancel-share').addEventListener('click', () => {
        document.getElementById('share-modal').classList.add('hidden');
        document.getElementById('share-form').reset();
        pendingShare = null;
    });

    document.getElementById('share-modal').addEventListener('click', (e) => {
        if (e.target.id === 'share-modal') {
            document.getElementById('share-modal').classList.add('hidden');
            document.getElementById('share-form').reset();
            pendingShare = null;
        }
    });

    document.getElementById('share-form').addEventListener('submit', (e) => {
        e.preventDefault();

        if (!pendingShare) return;

        const email = document.getElementById('share-email').value;
        const message = document.getElementById('share-message').value;
        const { versionNumber, versionMessage, content } = pendingShare;

        let emailBody = '';
        if (message) {
            emailBody += message + '\n\n';
        }
        emailBody += `File: ${currentFile.Name}\n`;
        emailBody += `Version: ${versionNumber} - ${versionMessage}\n\n`;
        emailBody += '--- Content ---\n\n';
        emailBody += content;

        const subject = `${currentFile.Name} - Version ${versionNumber}`;
        const mailtoLink = `mailto:${encodeURIComponent(email)}?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(emailBody)}`;

        window.location.href = mailtoLink;

        document.getElementById('share-modal').classList.add('hidden');
        document.getElementById('share-form').reset();
        pendingShare = null;
    });

    // Version content modal
    document.getElementById('close-modal-btn').addEventListener('click', () => {
        document.getElementById('version-modal').classList.add('hidden');
    });

    document.getElementById('version-modal').addEventListener('click', (e) => {
        if (e.target.id === 'version-modal') {
            document.getElementById('version-modal').classList.add('hidden');
        }
    });

    // Copy version content
    document.getElementById('copy-btn').addEventListener('click', async () => {
        const content = document.getElementById('version-content').textContent;

        if (!content || content === 'No content available for this version.') {
            await Dialog.Info({
                title: 'Copy',
                message: 'No content to copy.'
            });
            return;
        }

        navigator.clipboard.writeText(content).then(() => {
            const btn = document.getElementById('copy-btn');
            const originalHTML = btn.innerHTML;

            btn.innerHTML = `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>`;

            setTimeout(() => {
                btn.innerHTML = originalHTML;
            }, 2000);
        }).catch(async (err) => {
            console.error('Failed to copy:', err);
            await Dialog.Error({
                title: 'Copy failed',
                message: 'Failed to copy to clipboard.'
            });
        });
    });

    // Save version form
    const saveBtn = document.getElementById('save-version-btn');
    const saveForm = document.getElementById('save-version-form');
    const cancelBtn = document.getElementById('timeline-cancel-btn');
    const versionForm = document.getElementById('timeline-version-form');

    saveBtn.addEventListener('click', () => {
        saveForm.classList.toggle('hidden');
    });

    cancelBtn.addEventListener('click', () => {
        saveForm.classList.add('hidden');
        versionForm.reset();
    });

    versionForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const fileInput = document.getElementById('timeline-file-input');
        const message = document.getElementById('timeline-message-input').value;

        if (!fileInput.files || fileInput.files.length === 0) {
            await Dialog.Info({
                title: 'No file selected',
                message: 'Please select a file to upload.'
            });
            return;
        }

        const file = fileInput.files[0];

        const fileName = file.name.toLowerCase();
        if (!fileName.endsWith('.txt') && !fileName.endsWith('.md')) {
            await Dialog.Error({
                title: 'Unsupported file type',
                message: 'ZenTxt only supports .txt and .md files.\n\nPlease convert your document to plain text before uploading.'
            });
            return;
        }

        const reader = new FileReader();

        reader.onload = async function(e) {
            const content = e.target.result;

            try {
                await SaveVersion(currentFileId, currentFile.Path, 'You', message, content);

                // Close form and reset
                saveForm.classList.add('hidden');
                versionForm.reset();

                // Reload timeline
                loadTimeline(currentFileId);
            } catch (error) {
                console.error('Error:', error);
                await Dialog.Error({
                    title: 'Error',
                    message: 'Failed to save version. Please try again.'
                });
            }
        };

        reader.onerror = async function() {
            await Dialog.Error({
                title: 'Error',
                message: 'Failed to read file. Please try again.'
            });
        };

        reader.readAsText(file);
    });
}

// ======================
// INIT
// ======================

document.addEventListener('DOMContentLoaded', () => {
    initRouter();
    initHomeDragDrop();
    initHomeModal();
    initTimelineEventHandlers();
    initTimelineModals();
});
