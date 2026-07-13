<script type="ts">

    import { SaveSettings, GetSettings, StartUpdate } from './../../wailsjs/go/main/App'

    import { getContext, onMount } from 'svelte';
    import type { Writable } from 'svelte/store';
        
    interface ClientUpdate {
        isLatest: boolean;
        isLoaded: boolean;
        isInMatch: boolean;
        gameOpen: boolean;
        tokenTest: string;
    }

    const clientDataStore = getContext<Writable<ClientUpdate>>('client-update');
    let clientData: ClientUpdate;

    let unsubscribe: (() => void) | null = null;
    
    if (clientDataStore) {

        unsubscribe = clientDataStore.subscribe(data => {
            clientData = data;
        });

    }

    let settings: Record<string, boolean> = {};

    onMount( () => {
        
        // Check current settings

        GetSettings().then( ( data: Record<string, boolean>) => {

            if ( data == null ) {
                return
            }

            settings = data

        })

        return () => {
            if (unsubscribe) unsubscribe();
        };

    })

    function updateSettings( SettingName: string, value: boolean ) {

        if ( settings == null ) {
            settings = {}
        }

        settings[SettingName] = value
        
        SaveSettings(settings)

    }

    function getSetting( SettingName: string ): Boolean {

        if (settings == null) {
            return false
        }

        return settings[SettingName] ?? false
    }
    
    var isUpdating: boolean = false

    function beginUpdate() {

        isUpdating = true

        StartUpdate()

    }

</script>

<main>

    {#if !clientData.isLatest}

        <div class ="update-setting"
        on:click={ () => { beginUpdate() } }
        on:keyup={ () => { beginUpdate() } }
        >
            {#if !isUpdating}

                <setting_name>Controller is out of date!</setting_name>
                <alt_name>Click to update</alt_name>

            {/if}

            {#if isUpdating}

                <setting_name>Updating..</setting_name>

            {/if}
        </div>

    {/if}

    <div class="settings-grid">

        <div class ="setting">
            <setting_name>Always On Top</setting_name>
            <value 
                class="toggle"
                on:click={ () => { updateSettings("AlwaysOnTop", !getSetting("AlwaysOnTop")) } }
                on:keyup={ () => { updateSettings("AlwaysOnTop", !getSetting("AlwaysOnTop")) } }
            >   
                <div class="toggle_thumb { settings[ "AlwaysOnTop" ] ?? false ? "on" : "off" }"/>
            </value>
        </div>

        <div class ="setting">
            <setting_name>Send Anonymous Data</setting_name>
            <value 
                class="toggle"
                on:click={ () => { updateSettings("SendAnonymousData", !getSetting("SendAnonymousData")) } }
                on:keyup={ () => { updateSettings("SendAnonymousData", !getSetting("SendAnonymousData")) } }
            >
                <div class="toggle_thumb { settings[ "SendAnonymousData" ] ?? false ? "on" : "off" }"/>
            </value>
        </div>

        <!--<div class ="setting">
            <setting_name>Is Latest Version?</setting_name>
            <value 
                class="toggle"
            >   
                <div class="toggle_thumb { clientData.isLatest ? "on" : "off" }"/>
            </value>
        </div>-->

    </div>

</main>

<style>

    .settings-grid {
        
        display: grid;
        grid-template-columns: 48.75% 48.75%;
        grid-template-rows: repeat(auto-fit, minmax(2.5rem, max-content));

        gap: 0.2rem 2.5%;

        margin-left: 2.5%;
        margin-top: 0.2rem;
        background-color: transparent;

        position: absolute;

        flex-direction: column;

        width: 95%;
        height: 100%;

    }
    
    .update-setting {

        margin-top: 0.2rem;

        user-select: none;

        color: hsla(180, 67%, 99%, 0.7);
        font-weight: 700;
        font-family: 'DMSans', sans-serif;

        border-radius: 0.5rem;

        display: flex;
        flex-direction: column;

        margin-left: 2.5%;

        width: 95%;
        height: 2.5rem;

        background-color: hsl(49, 60%, 35%);

        justify-items: center;
        align-items: center;

        transition: 500ms all cubic-bezier(0.075, 0.82, 0.165, 1);

    }
    
    .update-setting:hover {

        background-color: hsl(49, 60%, 20%);

    }

    .update-setting alt_name {

        width: 100%;
        
        font-weight: 600 !important;

        color: rgb(105, 15, 15) !important;

        font-size: 0.7rem;

        transform: translateY(-20%);

    }

    .update-setting setting_name {

        height: 100%;
        width: 100%;

        align-content: center;
        justify-content: center;
        line-height: 100%;
    }

    .setting {

        user-select: none;

        color: hsla(180, 67%, 99%, 0.7);
        font-weight: 700;
        font-family: 'DMSans', sans-serif;

        margin-top: 0.2rem;

        border-radius: 0.5rem;

        display: flex;

        width: 100%;
        height: 2.5rem;

        background-color: hsl(0, 0%, 10%);

        justify-items: center;
        align-items: center;

    }

    .setting setting_name {

        height: 100%;
        width: calc( 100% - 5rem );

        align-content: center;
        justify-content: center;
        line-height: 100%;
    }

    .toggle {
        
        align-content: center;
        justify-content: center;
        
        background-color: hsl(0, 0%, 15%);
        border: 1.5px white solid;

        width: 3rem;
        height: 1.5rem;

        border-radius: 0.3rem;

    }

    .toggle_thumb {

        position: relative;

        width: 1.5rem;
        height: 1.5rem;

        transition: 250ms all;

    }

    .toggle .on {
        background-color: rgb(12, 189, 12);
        transform: translateX(100%);
        border-radius: 0 0.3rem 0.3rem 0;
    }

    .toggle .off {
        background-color: red;
        border-radius: 0.3rem 0 0 0.3rem;
    }

    .toggle::after {

        content: "";
        position: absolute;

        top: 0.7rem;

        width: 1px;
        height: 1.5rem;
        background-color: white;

    }

</style>