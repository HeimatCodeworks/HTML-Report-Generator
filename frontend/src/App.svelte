<script>
  import CompanySelect from "./lib/CompanySelect.svelte";
  import ReportSelect from "./lib/ReportSelect.svelte";
  import ReportView from "./lib/ReportView.svelte";
  import EmailSender from "./lib/EmailSender.svelte";

  let currentView = "home";
  let selectedCompany = null;
  let selectedReport = null;
  let generatedHtml = "";
  let showEmailModal = false;

  function handleCompanySelect(event) {
    selectedCompany = event.detail;
    currentView = "report";
  }

  function handleReportSelect(event) {
    selectedReport = event.detail;
    generatedHtml = event.detail.html;
    currentView = "result";
  }

  function startOver() {
    currentView = "company";
  }

  function startOverForCompany() {
    currentView = "report";
  }

  function handleGoBack() {
    currentView = "company";
  }

</script>

<main>
  {#if currentView === 'home'}
    <div class="home-view">
        <h1>Report Generator</h1>
        <button class="btn" on:click={() => currentView = 'company'}>Generate New Report</button>
    </div>
  {:else if currentView === 'company'}
    <CompanySelect on:companySelected={handleCompanySelect} />
  {:else if currentView === 'report'}
    <ReportSelect 
      company={selectedCompany} 
      on:reportGenerated={handleReportSelect} 
      on:goBack={handleGoBack}
    />
  {:else if currentView === 'result'}
    <ReportView 
      company={selectedCompany} 
      report={selectedReport} 
      bind:html={generatedHtml}
      on:startOver={startOver}
      on:startOverForCompany={startOverForCompany}
      on:sendEmail={() => showEmailModal = true}
    />
  {/if}
  {#if showEmailModal}
    <EmailSender
      company={selectedCompany}
      report={selectedReport}
      html={generatedHtml}
      on:close={() => showEmailModal = false}
    />
  {/if}
</main>

<style global>
  :root {
    --primary-dark-gray: #1f2937; /* Main background */
    --secondary-dark-gray: #374151; /* Panel backgrounds, input fields */
    --accent-blue: #3b82f6; /* Buttons, highlights */
    --accent-blue-hover: #2563eb; /* Button hover */
    --primary-text: #f3f4f6; /* Main text color */
    --secondary-text: #9ca3af; /* Lighter text for placeholders, etc. */
    --border-color: #4b5563;
  }

  main {
    background-color: var(--primary-dark-gray);
    height: 100vh;
    padding: 2rem;
    box-sizing: border-box;
  }

  h1 {
    font-weight: 500;
  }

  .btn {
    background-color: var(--accent-blue);
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.2s;
  }
  .btn:hover {
    background-color: var(--accent-blue-hover);
  }

</style>