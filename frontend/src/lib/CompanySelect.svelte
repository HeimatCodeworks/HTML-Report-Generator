<script>
    import { onMount, createEventDispatcher } from "svelte";
    import { GetCompanies } from "wailsjs/go/main/App";

    const dispatch = createEventDispatcher();

    let allCompanies = [];
    let filteredCompanies = [];
    let searchTerm = "";
    let isLoading = true;

    onMount(async () => {
        allCompanies = await GetCompanies();
        filteredCompanies = allCompanies;
        isLoading = false;
    });

    $: {
        if (searchTerm) {
            filteredCompanies = allCompanies.filter(c =>
                c.name.toLowerCase().includes(searchTerm.toLowerCase())
            );
        } else {
            filteredCompanies = allCompanies;
        }
    }

    function selectCompany(company) {
        dispatch('companySelected', company);
    }
</script>

<div class="container">
    <h2>1. Select a Company</h2>
    <input type="text" bind:value={searchTerm} placeholder="Search for a company..." class="search-input"/>

    {#if isLoading}
        <p>Loading companies...</p>
    {:else}
        <ul class="item-list">
            {#each filteredCompanies as company (company.teamId)}
                <li>
                    <button on:click={() => selectCompany(company)}>
                        {company.name}
                    </button>
                </li>
            {/each}
        </ul>
    {/if}
</div>

<style>
    .container h2 {
        border-bottom: 1px solid var(--border-color);
        padding-bottom: 0.5rem;
        margin-top: 0;
    }
    .search-input { 
        width: 100%; 
        padding: 12px; 
        margin-bottom: 1rem; 
        background-color: var(--secondary-dark-gray);
        border: 1px solid var(--border-color);
        border-radius: 6px;
        color: var(--primary-text);
        font-size: 1rem;
    }
    .search-input::placeholder {
        color: var(--secondary-text);
    }
    .item-list { 
        list-style: none; 
        padding: 0; 
        max-height: 60vh; 
        overflow-y: auto; 
    }
    .item-list li { 
        padding: 0; 
        border: 1px solid transparent; /* Hide border until hover */
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