<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { SendEmail, GetUsersForCompany } from "wailsjs/go/main/App";

    export let company;
    export let report;
    export let html;

    const dispatch = createEventDispatcher();

    let toRecipients = [];
    let ccRecipients = [];
    let toInput = "";
    let ccInput = "";

    onMount(async () => {
        try {
            const emails = await GetUsersForCompany(company.teamId);
            emails.forEach(email => addRecipient(email, 'to'));
        } catch (err) {
            statusMessage = `Could not fetch users: ${err}`;
        }
    });

    function handleKeydown(event, type) {
        if (event.key === 'Enter' || event.key === ',' || event.key === ' ') {
            event.preventDefault();
            addRecipient(type === 'to' ? toInput : ccInput, type);
            if (type === 'to') toInput = "";
            if (type === 'cc') ccInput = "";
        }
    }

    function addRecipient(email, type) {
        const trimmedEmail = email.trim();
        if (trimmedEmail && (type === 'to' ? !toRecipients.includes(trimmedEmail) : !ccRecipients.includes(trimmedEmail))) {
            if (type === 'to') {
                toRecipients = [...toRecipients, trimmedEmail];
            } else {
                ccRecipients = [...ccRecipients, trimmedEmail];
            }
        }
    }

    function removeRecipient(email, type) {
        if (type === 'to') {
            toRecipients = toRecipients.filter(e => e !== email);
        } else {
            ccRecipients = ccRecipients.filter(e => e !== email);
        }
    }


    let subject = `${company.name} - ${report.name} Report`;
    let statusMessage = "";
    let isSending = false;

    async function handleSendEmail() {
        isSending = true;
        statusMessage = "Sending...";
        try {
            await SendEmail(toRecipients, ccRecipients, subject, html);
            statusMessage = "Email sent successfully!";
            setTimeout(() => dispatch('close'), 2000);
        } catch (err) {
            statusMessage = `Error: ${err}`;
        } finally {
            isSending = false;
        }
    }
</script>

<div class="modal-backdrop">
    <div class="modal-content">
        <div class="modal-header">
            <h3>Send Report via Email</h3>
            <button class="close-btn" on:click={() => dispatch('close')}>&times;</button>
        </div>
        
        <div class="form-group">
            <label for="to">To:</label>
            <div class="recipient-input-wrapper">
                {#each toRecipients as email}
                    <span class="pill">
                        {email}
                        <button on:click={() => removeRecipient(email, 'to')}>&times;</button>
                    </span>
                {/each}
                <input id="to" type="text" bind:value={toInput} on:keydown={(e) => handleKeydown(e, 'to')} placeholder={toRecipients.length === 0 ? 'recipient@example.com' : ''}>
            </div>
        </div>
        <div class="form-group">
            <label for="cc">Cc:</label>
            <div class="recipient-input-wrapper">
                 {#each ccRecipients as email}
                    <span class="pill">
                        {email}
                        <button on:click={() => removeRecipient(email, 'cc')}>&times;</button>
                    </span>
                {/each}
                <input id="cc" type="text" bind:value={ccInput} on:keydown={(e) => handleKeydown(e, 'cc')} placeholder={ccRecipients.length === 0 ? 'cc@example.com' : ''}>
            </div>
        </div>
        <div class="form-group">
            <label for="subject">Subject:</label>
            <input id="subject" type="text" bind:value={subject}>
        </div>

        <div class="preview-pane">
            <iframe srcdoc={html} title="Email Preview"></iframe>
        </div>

        <div class="modal-footer">
            {#if statusMessage}
                <p class="status">{statusMessage}</p>
            {/if}
            <button class="btn" on:click={handleSendEmail} disabled={isSending}>
                {isSending ? 'Sending...' : 'Send Email'}
            </button>
        </div>
    </div>
</div>

<style>
    .modal-backdrop { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.7); display: flex; justify-content: center; align-items: center; z-index: 100; }
    .modal-content { background-color: var(--primary-dark-gray); border: 1px solid var(--border-color); border-radius: 8px; width: 90%; max-width: 800px; display: flex; flex-direction: column; max-height: 90vh; }
    .modal-header { display: flex; justify-content: space-between; align-items: center; padding: 1rem 1.5rem; border-bottom: 1px solid var(--border-color); }
    .modal-header h3 { margin: 0; }
    .close-btn { background: none; border: none; color: var(--primary-text); font-size: 2rem; cursor: pointer; }
    .form-group { padding: 0.5rem 1.5rem; display: flex; align-items: center; gap: 1rem; }
    .form-group label { flex-shrink: 0; flex-basis: 80px; text-align: right; }
    .recipient-input-wrapper {
        flex-grow: 1;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 0.5rem;
        padding: 8px;
        background-color: var(--secondary-dark-gray);
        border: 1px solid var(--border-color);
        border-radius: 4px;
    }
    .recipient-input-wrapper input {
        flex-grow: 1;
        border: none;
        background: none;
        padding: 0;
        color: var(--primary-text);
        min-width: 150px;
    }
    .recipient-input-wrapper input:focus {
        outline: none;
    }
    .pill {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        background-color: var(--accent-blue);
        color: white;
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 0.9rem;
    }
    .pill button {
        background: none;
        border: none;
        color: white;
        font-weight: bold;
        cursor: pointer;
        padding: 0;
        line-height: 1;
    }
    
    .form-group #subject {
        flex-grow: 1; padding: 8px; background-color: var(--secondary-dark-gray);
        border: 1px solid var(--border-color); border-radius: 4px; color: var(--primary-text);
    }
    .preview-pane { flex-grow: 1; padding: 1rem 1.5rem; height: 450px; }
    .preview-pane iframe { width: 100%; height: 100%; border: 1px solid var(--border-color); background-color: #fff; }
    .modal-footer { padding: 1rem 1.5rem; border-top: 1px solid var(--border-color); display: flex; justify-content: flex-end; align-items: center; gap: 1rem; }
    .status { margin: 0; color: var(--secondary-text); flex-grow: 1; }
</style>