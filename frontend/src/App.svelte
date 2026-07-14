<script lang="ts">
  import {CloseWindow, MinimizeWindow, SetWindowSize, UpdateCurrentClient } from './../wailsjs/go/main/App'
  
  import { 
    setWorkingArea, 
    sizeAnimation,
  } from './utils/animations';

  import { onMount, setContext } from 'svelte';
  import { writable } from 'svelte/store';

  import { EventsOn, EventsOff, BrowserOpenURL } from '../wailsjs/runtime';

  import SettingsIcon from "./components/icons/settings.svelte"
  import Endpoints from './components/icons/endpoints.svelte';
  import History from './components/icons/history.svelte';
  import LoadoutsIcon from './components/icons/loadouts.svelte';

  import Loading from './components/loading.svelte';
  import Settings from './components/settings.svelte';
  import Matches from './components/matches.svelte';
  import Loadouts from './components/loadouts.svelte';
  import Randomizer from './components/randomizer.svelte';
    
  interface ClientUpdate {
    isLoaded: boolean;
	  isLatest: boolean;
    isInMatch: boolean;
    gameOpen: boolean;
    tokenTest: string;
  }

  let clientUpdate: ClientUpdate = {
    isLoaded: false,
    isLatest: true,
    isInMatch: false,
    gameOpen: false,
    tokenTest: ""
  }

  interface WindowSize {
    width: number;
    height: number;
  }

  interface RGB {
    r: number;
    g: number;
    b: number;
  }

  let clientUpdateStore = writable<ClientUpdate>(clientUpdate)
  setContext( "client-update", clientUpdateStore)

  function rgbToCss(color: RGB): string {
    return `${color.r}, ${color.g}, ${color.b}`;
  }

  interface Menu {
    layout: typeof Loading;
    windowSize: WindowSize;
    svg: typeof SettingsIcon;
    color: RGB;
  }

  interface Menus {
    loading: Menu;
    settings: Menu;
    matchs: Menu;
    loadouts: Menu;
    randomizer: Menu;
    endpoints: Menu;
  }

  const Menus: Menus = {
    loading: {
      layout: Loading,
      windowSize: {
        width: 650,
        height: 175
      },
      svg: SettingsIcon,
      color: {r: 0, g: 0, b: 0 },
    },
    settings: {
      layout: Settings,
      windowSize: {
        width: 650,
        height: 400
      },
      svg: SettingsIcon,
      color: {r: 255, g: 255, b: 255 },
    },
    matchs: {
      layout: Matches,
      windowSize: {
        width: 650,
        height: 38
      },
      svg: SettingsIcon,
      color: {r: 83, g: 135, b: 219 },
    },
    loadouts: {
      layout: Loadouts,
      windowSize: {
        width: 600,
        height: 400
      },
      svg: LoadoutsIcon,
      color: {r: 201, g: 114, b: 46 },
    },
    randomizer: {
      layout: Randomizer,
      windowSize: {
        width: 650,
        height: 400
      },
      svg: History,
      color: {r: 212, g: 68, b: 159 },
    },
    endpoints: {
      layout: Settings,
      windowSize: {
        width: 650,
        height: 400
      },
      svg: Endpoints,
      color: {r: 74, g: 17, b: 217 },
    },

    // Maybe Friend info?
    // Maybe Shop Info?
  };

  let CurrentActiveMenu = Menus.loading;

  let workingarea = document.body;

  function changeLayout( menu: Menu ) {

    setWorkingArea(workingarea)

    if (menu == Menus.matchs) {
      
      CurrentActiveMenu = menu;
      document.documentElement.style.setProperty('--menu-color', rgbToCss(menu.color));
      return

    }

    sizeAnimation(
      CurrentActiveMenu.windowSize.width, 
      CurrentActiveMenu.windowSize.height,
      menu.windowSize.width,
      menu.windowSize.height
    )

    CurrentActiveMenu = menu;

    document.documentElement.style.setProperty('--menu-color', rgbToCss(menu.color));

  }

  function setupMounts() {

    EventsOn( "updateClient", ( data: ClientUpdate) => {

      // Updates when server asks client to update

      clientUpdateStore.set(data)

      if (clientUpdate.isLoaded != data.isLoaded ) {

        if (data.isLoaded) {
          
          changeLayout( Menus.matchs )

        } else {
          
          changeLayout( Menus.loading )
          
        }

      }

      clientUpdate = data
      //console.log( "Update:", data )

    })

    UpdateCurrentClient()

  }

  setupMounts()
  
  onMount( () => {

    changeLayout( CurrentActiveMenu )
    
    return () => {
      console.log('Unmounting component - removing listeners');

      EventsOff('updateClient');
    };

  })

</script>

<main>

  <div class="working-area"
    bind:this={workingarea}
  >

    <div class="bar">

      <div class="icon-tray flex-left">
      
        <!-- Settings Button -->
        <button class="normal-button settings-button {!clientUpdate.isLoaded ? 'tray-inactive' : ''} {CurrentActiveMenu == Menus.settings ? 'active' : ''}" on:click={() => changeLayout(Menus.settings)}>
          <div class="settings-svg"><svelte:component this={Menus.settings.svg} /></div>
        </button>
      
        <!-- Matches Button -->
        <button class="normal-button {clientUpdate.isLoaded ? '' : 'tray-inactive'} {CurrentActiveMenu == Menus.matchs ? 'active' : ''}" on:click={() => changeLayout(Menus.matchs)}>
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width={CurrentActiveMenu == Menus.matchs ? '4' : '2'}
              stroke-linecap="round"
              stroke-linejoin="round"
              aria-hidden="true"
              >
              <circle cx="12" cy="12" r="8" />
            </svg>
        </button>
      
        <!-- Loadouts Button -->
        <button class="normal-button {clientUpdate.isLoaded ? '' : 'tray-inactive'} {CurrentActiveMenu == Menus.loadouts ? 'active' : ''}" on:click={() => changeLayout(Menus.loadouts)}>
          <div class="settings-svg"><svelte:component this={Menus.loadouts.svg} /></div>
        </button>
      
        <!-- History Button -->
        <button class="normal-button {clientUpdate.isLoaded ? '' : 'tray-inactive'} {CurrentActiveMenu == Menus.randomizer ? 'active' : ''}" on:click={() => changeLayout(Menus.randomizer)}>
          <div class="settings-svg"><svelte:component this={Menus.randomizer.svg} /></div>
        </button>
      
        <!-- Endpoints Button -->
        <button class="normal-button {clientUpdate.isLoaded ? '' : 'tray-inactive'} {CurrentActiveMenu == Menus.endpoints ? 'active' : ''}" on:click={() => changeLayout(Menus.endpoints)}>
          <div><svelte:component this={Menus.endpoints.svg} /></div>
        </button>

      </div>

      <div class="title">Valorant Controller</div>

      <div class="icon-tray flex-right">

        <!-- Exit Button -->
        <button class="danger-button" on:click={CloseWindow}>
          <svg width="10" height="10" viewBox="0 0 10 10" fill="none">
            <path d="M0.5 0.5L9.5 9.5M9.5 0.5L0.5 9.5" stroke="currentColor" style="stroke-width:2"></path>
          </svg>
        </button>
      
        <!-- Minimize Button -->
        <button class="normal-button" on:click={MinimizeWindow}>
          <svg width="10" height="2" viewBox="0 0 10 2" fill="none">
            <rect width="10" height="2" fill="currentColor" />
          </svg>
        </button>

      </div>

    </div>

    <!-- Main Window Area -->

    <div class="main-window">
      <svelte:component this={CurrentActiveMenu.layout}/>
    </div>

  </div>

</main>

<style>

  .working-area {
    position: absolute;
    display: flex;
    flex-direction: column;

    text-align: center;
    border-radius: 10px;
    background-color: rgb(12, 12, 14);
    outline: 1px solid rgb(42, 41, 56);
    outline-offset: -1px;

    overflow: hidden;
  }

  .bar {
    padding-inline: 0.75rem;
    box-sizing: border-box;
    --wails-draggable: drag;
    height: 2.5rem;
    background-color: rgb( 26, 26, 26 );
    border-top-left-radius: 10px;
    border-top-right-radius: 10px;
    border-bottom: 2px solid rgb( 48, 48, 48 );

    display: flex;
    align-items: center;
    justify-content: stretch;
  }

  .main-window {

    flex: 1;
    height: max-content;

    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;

  }

  .settings-svg {
    transition: transform 0.5s ease-in-out !important;
    transform: rotate(0deg);
  }

  .settings-button:hover .settings-svg {
    transform: rotate(240deg);
  }

  .icon-tray {
    flex: 1;
    height: 2rem;

    display: flex;
    gap: 0.25rem;

    transition: 0.5s all;
  }

  .tray-inactive {

    pointer-events: none;
    opacity: 0;

  }

  .flex-left { flex-direction: row; }
  .flex-right { flex-direction: row-reverse; }

  .title {
    user-select: none;

    color: hsla(180, 67%, 99%, 0.75);
    font-family: 'DMSans', sans-serif;
    font-weight: 800;
    font-size: 1rem;
    letter-spacing: -0.04em;
  }

  .icon-tray button {
    display: flex;
    border-radius: 0.375rem;
    border: 0;
    background-color: transparent;
    width: 2rem;
    height: 2rem;

    color: hsla(180, 67%, 99%, 0.5);

    align-items: center;
    justify-content: center;
    --wails-draggable: no-drag;

    cursor: pointer;

    transition: all 0.15s ease;
  }

  .active {

    color: rgb(var(--menu-color)) !important;
    background-color: rgba(var(--menu-color), 0.15) !important;

  }

  .danger-button:hover {
    background: rgba(248, 113, 113, 0.15);
    color: rgba(248, 113, 113);
  }

  .normal-button svg { transition: all 0.15s ease; }
  .normal-button div { transition: all 0.15s ease; }

  .normal-button div {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .normal-button:hover {
    background: hsl(0, 0%, 15%);
    color:  hsla(180, 67%, 99%, 0.9);
  }

  .icon-tray button svg {
    display: block;
  }

</style>
