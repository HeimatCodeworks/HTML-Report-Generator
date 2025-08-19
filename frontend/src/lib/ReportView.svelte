<script>
    import { createEventDispatcher, beforeUpdate } from "svelte";
    export let company;
    export let report;
    export let html;

    const dispatch = createEventDispatcher();
    let isEditing = false;
    let editTooltipText = 'Edit HTML';
    let copyHtmlTooltipText = 'Copy HTML';
    let copyPromptTooltipText = 'Copy AI Prompt';

    async function copyAiPrompt() {
        const masterPrompt = `
You are a surgical text-replacement tool with an expert-level understanding of digital marketing analysis. Your sole purpose is to find specific placeholder text within a provided block of HTML and replace it with a data-driven analysis. The data below is data from a company's Apple Search Ads campaigns and apps. This report will be sent as an email to clients who use our ASA automation tool. Try to spin the recommendations in the best way possible for us and our work. You can acknowledge need for improvement if it is obvious that it is needed.

**CRITICAL INSTRUCTIONS:**
1.  **FIND:** Scan the HTML code below and find every exact instance of the placeholder text: "INSERT AI SUMMARY HERE".
2.  **ANALYZE:** For each placeholder, analyze the relevant data tables within the HTML to generate a concise, professional summary.
    - Make insightful observations and recommendations in response to the data.
    - The first placeholder is for the "Executive Summary" and should analyze the "Overall Account Performance" table.
    - Subsequent placeholders are for individual campaigns and should analyze that specific campaign's performance table.
3.  **REPLACE:** Replace ONLY the placeholder text "INSERT AI SUMMARY HERE" with your generated analysis. The campaign analysis must begin with a bolded "Analysis:" tag (&lt;strong&gt;Analysis:&lt;/strong&gt;).
4.  **PRESERVE FORMATTING:** You MUST NOT alter any whitespace, indentation, newlines, or any other formatting in the original HTML. The structure must be preserved exactly as provided. Your output should be the original HTML with only the placeholder text surgically replaced. Do not insert any "returns" or anything that may change the spacing of the html text blocks.
5.  **OUTPUT HTML ONLY:** Your entire response must be ONLY the complete, valid HTML code. Do not include any conversational text, explanations, or markdown code fences (like \`\`\`html) in your response. Put the code in an easily copyable code block as the html will be copied and pasted directly.

Here is the HTML report:

${html}
`;
        try {
            await navigator.clipboard.writeText(masterPrompt.trim());
            copyPromptTooltipText = "Copied!";
            setTimeout(() => copyPromptTooltipText = "Copy AI Prompt", 2000);
        } catch (err) {
            alert('Failed to copy AI Prompt.');
        }
    }

    async function copyHtml() {
        try {
            await navigator.clipboard.writeText(html);
            copyHtmlTooltipText = "Copied!";
            setTimeout(() => copyHtmlTooltipText = "Copy HTML", 2000);
        } catch (err) {
            alert('Failed to copy HTML.');
        }
    }

    function toggleEdit() {
        isEditing = !isEditing;
        editTooltipText = isEditing ? 'Lock HTML' : 'Edit HTML';
    }

    // --- Live Preview Logic ---
    let iframeElement;
    let savedScrollY = 0;
    let debounceTimer;
    let previewHtml = html;
    $: {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            previewHtml = html;
        }, 250);
    }
    beforeUpdate(() => {
        if (iframeElement && iframeElement.contentWindow) {
            savedScrollY = iframeElement.contentWindow.scrollY;
        }
    });
    function restoreScrollPosition() {
        if (iframeElement && iframeElement.contentWindow) {
            iframeElement.contentWindow.scrollTo(0, savedScrollY);
        }
    }
</script>

<div class="container">
    <h3>{report.name} for {company.name}</h3>
    <div class="split-view">
        <div class="panel">
            <h4>HTML Source</h4>
            <div class="html-view">
                <div class="inlaid-buttons">
                    <button class="copy-btn" on:click={toggleEdit} title={editTooltipText}>
                        {#if isEditing}
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="18" height="18"><path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zM9 8V6c0-1.66 1.34-3 3-3s3 1.34 3 3v2H9z"/></svg>
                        {:else}
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="18" height="18"><path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/></svg>
                        {/if}
                    </button>
                    <button class="copy-btn" on:click={copyHtml} title={copyHtmlTooltipText}>
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="18" height="18"><path d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/></svg>
                        <span class="tooltip">{copyHtmlTooltipText}</span>
                    </button>
                    <button class="copy-btn" on:click={copyAiPrompt} title={copyPromptTooltipText}>
                         <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="18" height="18">
                            <path d="M15 4H5v16h14V8h-4V4zM3 2.992C3 2.444 3.447 2 3.999 2H16l5 5v13.993A1.007 1.007 0 0 1 20.007 22H3.993A1.007 1.007 0 0 1 3 21.008V2.992zM11 14v-2h2v2h2v2h-2v2h-2v-2H9v-2h2z"/>
                        </svg>
                        <span class="tooltip">{copyPromptTooltipText}</span>
                    </button>
                </div>
                <textarea class="html-textarea" bind:value={html} disabled={!isEditing} spellcheck="false"></textarea>
            </div>
        </div>
        <div class="panel">
            <h4>Live Preview</h4>
            <div class="preview-view">
                <iframe bind:this={iframeElement} srcdoc={previewHtml} title="Report Preview" frameborder="0" on:load={restoreScrollPosition}></iframe>
            </div>
        </div>
    </div>
    <div class="button-group">
        <button class="btn-secondary" on:click={() => dispatch('sendEmail')}>Send as Email</button>
        <button class="btn-secondary" on:click={() => dispatch('startOverForCompany')}>Generate Another Report for {company.name}</button>
        <button class="btn-secondary" on:click={() => dispatch('startOver')}>Generate New Report</button>
    </div>
</div>

<style>
    .container h3 { border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem; margin-top: 0; }
    .split-view { display: flex; gap: 1rem; width: 100%; margin-bottom: 1rem; }
    .panel { flex: 1; display: flex; flex-direction: column; min-width: 0; }
    .panel h4 { margin-top: 0; margin-bottom: 0.5rem; color: var(--secondary-text); font-weight: 500; }
    .html-view, .preview-view { background-color: var(--secondary-dark-gray); border: 1px solid var(--border-color); border-radius: 6px; height: 60vh; overflow: hidden; }
    .html-view { position: relative; padding: 0; }
    .inlaid-buttons { position: absolute; top: 8px; right: 8px; display: flex; gap: 6px; z-index: 10; }
    .copy-btn { position: relative; background-color: #4b5563; border: 1px solid var(--border-color); color: var(--primary-text); padding: 6px; border-radius: 6px; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: background-color 0.2s; }
    .copy-btn:hover { background-color: var(--accent-blue); }
    .copy-btn .tooltip { display: none; position: absolute; bottom: -30px; right: 0; background-color: #111827; color: white; padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; white-space: nowrap; }
    .copy-btn:hover .tooltip { display: block; }
    .html-textarea { width: 100%; height: 100%; margin: 0; padding: 1rem; box-sizing: border-box; background-color: transparent; border: none; color: var(--primary-text); font-family: 'Courier New', Courier, monospace; font-size: 0.9rem; resize: none; white-space: pre; overflow-wrap: normal; overflow-x: auto; }
    .html-textarea:focus { outline: none; }
    .html-textarea:disabled { color: #9ca3af; }
    .preview-view iframe { width: 100%; height: 100%; background-color: #fff; border-radius: 6px; }
    .button-group { margin-top: 1rem; display: flex; justify-content: center; gap: 1rem; }
    .btn-secondary { background-color: var(--secondary-dark-gray); color: var(--primary-text); border: 1px solid var(--border-color); padding: 10px 20px; border-radius: 6px; cursor: pointer; font-size: 1rem; transition: border-color 0.2s; }
    .btn-secondary:hover { border-color: var(--accent-blue); }
</style>