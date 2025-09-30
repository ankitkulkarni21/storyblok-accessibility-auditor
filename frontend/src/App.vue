<template>
  <div id="app">
    <div class="particles-container">
      <Particles
        :particle-count="200"
        :particle-spread="10"
        :speed="0.5"
        :particle-colors="['#f0ead6']"
        :move-particles-on-hover="true"
        :particle-hover-factor="0.1"
        :alpha-particles="false"
        :particle-base-size="100"
        :size-randomness="1"
        :camera-distance="20"
        :disable-rotation="false"
        class="w-full h-full"
      />
    </div>

    <div class="content-overlay">
      <a href="https://github.com/ankitkulkarni21/" target="_blank" class="github-link">
        <svg viewBox="0 0 16 16" version="1.1" width="32" height="32" aria-hidden="true" fill="white">
          <path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path>
        </svg>
      </a>
      <a href="https://www.linkedin.com/in/ankitkulkarni21/" target="_blank" class="linkedin-link">
        <svg xmlns="http://www.w3.org/2000/svg" width="34" height="34" fill="currentColor" class="bi bi-linkedin" viewBox="0 0 16 16">
          <path d="M0 1.146C0 .513.526 0 1.175 0h13.65C15.474 0 16 .513 16 1.146v13.708c0 .633-.526 1.146-1.175 1.146H1.175C.526 16 0 15.487 0 14.854zm4.943 12.248V6.169H2.542v7.225zm-1.2-8.212c.837 0 1.358-.554 1.358-1.248-.015-.709-.52-1.248-1.342-1.248S2.4 3.226 2.4 3.934c0 .694.521 1.248 1.327 1.248zm4.908 8.212V9.359c0-.216.016-.432.08-.586.173-.431.568-.878 1.232-.878.869 0 1.216.662 1.216 1.634v3.865h2.401V9.25c0-2.22-1.184-3.252-2.764-3.252-1.274 0-1.845.7-2.165 1.193v.025h-.016l.016-.025V6.169h-2.4c.03.678 0 7.225 0 7.225z"/>
        </svg>
      </a>
      <h1>Storyblok Accessibility Auditor</h1>
      <div class="centered-wrapper">

        <p class="subtitle">Enter a Story Slug and your Public Storyblok Token below to run an accessibility audit on your space.</p>       
        
        <div class="rule-selection">
          <label v-for="rule in availableRules" :key="rule.id" class="rule-toggle">
            <input type="checkbox" v-model="rule.active" :disabled="isLoading">
            {{ rule.label }}
          </label>
        </div>
        
        <form @submit.prevent="fetchAndAudit" class="input-form">
          <input 
            type="text" 
            v-model="userToken" 
            placeholder="Enter public Storyblok token (optional)" 
            class="token-input"
          /> 
          <input 
            type="text" 
            v-model="storySlug" 
            placeholder="Enter Story Slug (e.g., 'audit-test')" 
            class="slug-input"
          /> 
          <button type="submit" :disabled="isLoading" class="audit-button">
            {{ isLoading ? 'Auditing...' : 'Run Audit' }}
          </button>

          <button type="button" @click="clearResults" :disabled="isLoading" class="reset-button">
            Reset
          </button>
        </form>
        <div class="app-description">
            <p>This is the "Storyblok Accessibility Auditor", a full-stack tool built with Vue.js and a Go API backend.</p>
            <p>The auditor dynamically checks your Storyblok content for common structural accessibility issues based on the rules selected above. Currently checking for:</p>
            <ul>
                <li>H1 Count : Ensures only one "h1" heading tag exists per page for proper hierarchy.</li>
                <li>Placeholder Text : Scans for "lorem ipsum" placeholder text left by content editors.</li>
                <li>Image Alt Text : Verifies that all Image components have non-empty "alt" text.</li>
            </ul>
            <p>Results provide immediate visual feedback showing the exact content block where the issue was found.</p>
        </div>
        <div v-if="isLoading" class="loading-message">
          <div class="spinner"></div>
          Loading audit report (Connecting to the backend hosted on render)...
        </div>

        <div v-if="error" class="error-message">
          <span class="icon">❌</span> {{ error }}
        </div>

        <div v-if="auditReport" class="audit-results">
          <div :class="auditReport.Status === 'Pass' ? 'status-pass' : 'status-fail'">
            Audit Status: {{ auditReport.Status }}
            <span class="icon">{{ auditReport.Status === 'Pass' ? '✅' : '🚨' }}</span>
            <span class="issue-count">
              Total Issues: {{ auditReport.Issues ? auditReport.Issues.length : 0 }}
            </span>
          </div>
          
          <div v-if="auditReport.Issues && auditReport.Issues.length">
            <h2>Detected Issues:</h2>
            <ul class="issues-list">
              <li v-for="(issue, index) in auditReport.Issues" :key="index">
                  <div>
                      <strong>{{ issue.Component }}:</strong> {{ issue.Message }}
                      <span :class="'severity-' + issue.Severity.toLowerCase()">{{ issue.Severity }}</span>
                  </div>
                  
                  <pre v-if="issue.ProblemData" class="problem-code">
                      <code>{{ JSON.stringify(issue.ProblemData, null, 2) }}</code>
                  </pre>
              </li>
            </ul>
          </div>
          
          <div v-else class="no-issues">
            <span class="icon">🎉</span> No issues found! Your content looks good.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import Particles from "./Particles.vue";

const storySlug = ref('');
const userToken = ref('');
const auditReport = ref(null);
const isLoading = ref(false);
const error = ref(null);

const clearResults = () => {
    storySlug.value = '';
    userToken.value = '';
    auditReport.value = null;
    error.value = null;
};

const availableRules = ref([
    { id: 'h1_count', label: 'Check H1 Count (Max 1)', active: true },
    { id: 'lorem_ipsum', label: 'Check for Placeholder Text', active: true },
    { id: 'alt_text', label: 'Check Image Alt Text', active: true },
]);

const fetchAndAudit = async () => {
    isLoading.value = true;
    error.value = null;
    auditReport.value = null;

    try {
        const tokenToUse = userToken.value || import.meta.env.VITE_STORYBLOK_TOKEN; 

         if (!tokenToUse) { 
            error.value = "A Storyblok token is required. Please enter a public token or ensure VITE_STORYBLOK_TOKEN is set."; 
            isLoading.value = false; 
            return; 
        } 

        const rulesToSend = availableRules.value
            .filter(rule => rule.active)
            .map(rule => rule.id); 

        const sbResponse = await axios.get( 
            `https://api.storyblok.com/v2/cdn/stories/${storySlug.value}`, 
             {
             params: {
                  token: tokenToUse, 
                  version: 'draft',
                }
             }
        );

        const content = sbResponse.data.story.content; 

        const backendURL = 'https://storyblok-accessibility-auditor.onrender.com/api/audit';

        const response = await axios.post(backendURL, { 
            content: JSON.stringify(content),
            activeRules: rulesToSend,
        });

        auditReport.value = response.data;

    } catch (err) {
        console.error("Audit failed:", err);

        if (err.response) {
            if (err.response.status === 404) { 
                error.value = `Story not found for slug: "${storySlug.value}". Check the slug and ensure it exists in the space associated with the provided token.`; 
            } else if (err.response.status === 403) { 
                 error.value = "Storyblok Token Rejected (403 Forbidden). Please verify your Public Access Token is correct."; 
            } else if (err.response.status === 401) { 
                 error.value = "Storyblok Token Missing or Invalid (401). Check token access/validity."; 
            } else { 
                 error.value = `Error from Storyblok/Backend: ${err.response.status} - ${err.response.data.message || 'Check network or server logs.'}`; 
            }
           } else if (err.code === 'ERR_NETWORK') { 
            error.value = "Network Error: Could not connect to backend. Ensure Go server is running and the Render URL is correct.";
           } else { 
             error.value = "An unknown error occurred during the audit process."; 
        }
     } finally {
        isLoading.value = false;
    }
};
</script>

<style scoped>
.github-link {
	position: absolute;
	top: 70px;
	right: 40px;
	z-index: 20;
	opacity: 0.8;
	transition: opacity 0.3s ease;
}

.github-link:hover {
	opacity: 1;
}

.github-link:visited {
	color: #ffffff;
}

.linkedin-link {
	position: absolute;
	top: 70px;
	right: 90px;
	z-index: 20;
	opacity: 0.8;
	transition: opacity 0.3s ease;
}

.linkedin-link:hover {
	opacity: 1;
}

.linkedin-link:visited {
	color: #ffffff;
}

.linkedin-link svg {
	fill: white;
}
</style>