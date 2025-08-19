<script>
    import { onMount, createEventDispatcher } from "svelte";
    import { GetReports, GenerateReport } from "wailsjs/go/main/App";

    export let company;
    const dispatch = createEventDispatcher();

    let reports = [];
    let isLoading = false;

    onMount(async () => {
        reports = await GetReports();
    });

    async function selectReport(report) {
        isLoading = true;
        try {
            const html = await GenerateReport(company.teamId, report.id);
            dispatch('reportGenerated', { ...report, html });
        } catch (err) {
            alert("Error generating report: " + err);
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="container">
    <div class="title-header">
        <button class="btn-back" on:click={() => dispatch('goBack')} title="Go Back">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="20" height="20">
                <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z"/>
            </svg>
        </button>
        <h2>2. Select a Report for {company.name}</h2>
    </div>

    {#if isLoading}
        <p>Generating report, please wait...</p>
    {:else}
        <ul class="item-list">
            {#each reports as report (report.id)}
                 <li>
                    <button on:click={() => selectReport(report)}>
                        {report.name}
                    </button>
                </li>
            {/each}
        </ul>
    {/if}
</div>

<style>
    .title-header {
        display: flex;
        align-items: center;
        gap: 1rem; /* Space between button and title */
        border-bottom: 1px solid var(--border-color);
        padding-bottom: 0.5rem;
    }

    .title-header h2 {
        border-bottom: none;
        padding-bottom: 0;
        margin: 0;
    }

    .btn-back {
        background-color: var(--secondary-dark-gray);
        color: var(--primary-text);
        border: 1px solid var(--border-color);
        padding: 8px;
        border-radius: 6px;
        cursor: pointer;
        font-size: 1rem;
        transition: border-color 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .btn-back:hover {
        border-color: var(--accent-blue);
    }

    .item-list { 
        list-style: none; 
        padding: 0; 
        margin-top: 1rem;
        max-height: 60vh; 
        overflow-y: auto; 
    }
    .item-list li { 
        padding: 0; 
        border: 1px solid transparent;
        border-radius: 6px; 
        margin-bottom: 8px; 
        transition: all 0.2s;
    }
    .item-list li:hover { 
        border-color: var(--accent-blue); 
    }
    .item-list li button {
        background: var(--secondary-dark-gray);
        border: none;
        color: var(--primary-text);
        font-size: 1rem;
        font: inherit;
        padding: 16px;
        width: 100%;
        text-align: left;
        cursor: pointer;
        border-radius: 6px;
    }
</style>