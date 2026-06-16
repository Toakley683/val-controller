<script lang="ts">

    import { getContext, onMount } from 'svelte';
    import type { Writable } from 'svelte/store';
        
    interface ClientUpdate {
        isLoaded: boolean;
        isInMatch: boolean;
        gameOpen: boolean;
        tokenTest: string;
    }

    const clientDataStore = getContext<Writable<ClientUpdate>>('client-update');
    
    let clientData: ClientUpdate;

    let text = ""
    let unsubscribe: (() => void) | null = null;
    
    if (clientDataStore) {

        unsubscribe = clientDataStore.subscribe(data => {
        clientData = data;
        console.log('Client data updated:', data);
        });

        if (clientData.gameOpen == false ) {
            text = "Please open the game to continue"
        }

    }

    onMount( () => {

        console.log("ClientData")
        console.log(clientData)

        return () => {
            if (unsubscribe) unsubscribe();
        };

    })

</script>

<main>

    <div class="loader">
        <svg 
            width=32
            height=32
            fill="#ffffff" 
            viewBox="0 0 32 32" 
            role="img" 
            xmlns="http://www.w3.org/2000/svg" 
            data-darkreader-inline-fill="" 
            style="--darkreader-inline-fill: var(--darkreader-background-000000, #000000);">
            <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
            <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
            <g id="SVGRepo_iconCarrier">
            <title>Riot Games icon</title>
            <path d="M12.534 21.77l-1.09-2.81 10.52.54-.451 4.5zM15.06 0L.307 6.969 2.59 17.471H5.6l-.52-7.512.461-.144 1.81 7.656h3.126l-.116-9.15.462-.144 1.582 9.294h3.31l.78-11.053.462-.144.82 11.197h4.376l1.54-15.37Z"></path></g></svg></div>

    <p class="error_text">{text}</p>

</main>

<style>

    .error_text {
        position: absolute;
        bottom: -8px;
        align-self:flex-end;
    }

    main {
        width: 100%;
        height: 100%;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .loader {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 72px;
        height: 72px;
    }
    .loader:before,
    .loader:after {
        content: "";
        position: absolute;
        border-radius: 50%;
        animation: pulsOut 1.8s ease-in-out infinite;
        filter: drop-shadow(0 0 0.5rem rgba(255, 255, 255, 0.75));
    }
    .loader:before {
        width: 100%;
        padding-bottom: 100%;
        box-shadow: inset 0 0 0 0.5rem #fff;
        animation-name: pulsIn;
    }
    .loader:after {
        width: calc(100% - 1rem);
        padding-bottom: calc(100% - 1rem);
        box-shadow: 0 0 0 0 #fff;
    }

    @keyframes pulsIn {
        0% {
        box-shadow: inset 0 0 0 0.5rem #fff;
        opacity: 1;
        }
        50%, 100% {
        box-shadow: inset 0 0 0 0 #fff;
        opacity: 0;
        }
    }

    @keyframes pulsOut {
        0%, 50% {
        box-shadow: 0 0 0 0 #fff;
        opacity: 0;
        }
        100% {
        box-shadow: 0 0 0 0.5rem #fff;
        opacity: 1;
        }
    }

    .loader svg {
        position: absolute;
        top: 24px;
        left: 24px;

        animation: 1s rotate ease infinite;
    }
    
    @keyframes rotate {
        0%{    transform: rotate(0deg)}
        25%{    transform: rotate(-5deg)}
        50%{    transform: rotate(5deg)}
        75%{    transform: rotate(0)}
    }
      

</style>