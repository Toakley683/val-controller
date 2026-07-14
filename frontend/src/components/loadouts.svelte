<script type="ts">

    import { onMount } from 'svelte';
    import { EventsOn, EventsOff } from '../../wailsjs/runtime';
    import { valorantapi, main } from '../../wailsjs/go/models';
    import { SaveCurrentLoadout, DeleteSavedLoadout, LoadSavedLoadout, GetLoadouts, AddWeaponsToBeRandomized } from '../../wailsjs/go/main/App';
    
    import { 
        sizeAnimation,
        getWorkingAreaSize
    } from './../utils/animations';
    import { element } from 'svelte/internal';

    let transferNode: HTMLElement;
    let isLoadoutShown: boolean = false;

    let loadoutSelected: LoadoutItem

    const BuddyOffsets = {
        "63e6c2b6-4a8e-869c-3d4c-e38355226584": 70,
        "55d8a0f4-4274-ca67-fe2c-06ab45efdf58": 50,
        "9c82e19d-4575-0200-1a81-3eacf00cf872": 25,
        "ae3de142-4d85-2547-dd26-4e90bed35cf7": 35,
        "ee8e8d15-496b-07ac-e5f6-8fae5d4c7b1a": 45,
        "ec845bf4-4f79-ddda-a3da-0db3774b2794": 30,
        "910be174-449b-c412-ab22-d0873436b21b": 70,
        "44d4e95c-4157-0037-81b2-17841bf2e8e3": -5,
        "29a0cfab-485b-f5d5-779a-b59f85e204a8": -5,
        "1baa85b4-4c70-1284-64bb-6481dfc3bb4e": 10,
        "e336c6b8-418d-9340-d77f-7a9e4cfe0702": -5,
        "42da8ccc-40d5-affc-beec-15aa47b42eda": 25,
        "a03b24d3-4319-996d-0f8c-94bbfba1dfc7": 75,
        "4ade7faa-4cf1-8376-95ef-39884480959b": 65,
        "c4883e50-4494-202c-3ec3-6b8a9284f00b": 85,
        "462080d1-4035-2937-7c09-27aa2a5c27a7": 25,
        "f7e1b454-4ad4-1063-ec0a-159e56b58941": 5,
        "2f59173c-4bed-b6c3-2191-dea9b58be9c7": 0,
        "5f0aaf7a-4289-3998-d5ff-eb9a5cf7ef5c": 85,
        "410b2e0b-4ceb-1321-1727-20858f7f3477": 5,
    };
    
    function showLoadout( isShown: boolean, selectedLoadout: LoadoutItem ) {
        
        error = ""

        if (transferNode == null) {
            return
        }

        isLoadoutShown = isShown

        if (isShown ) {
            transferNode.style.transform = "translate(-50%, 0px)";
        } else {
            transferNode.style.transform = "translate(0%, 0px)";
        }

        resize()

        if (selectedLoadout != null) {
            loadoutSelected = selectedLoadout
        } else {
            loadoutSelected = null
        }

    }
    
    function resize() {

        if (transferNode != null ) {
            if ( isLoadoutShown ) {

                const size = getWorkingAreaSize()

                sizeAnimation(
                    size.x,
                    size.y,
                    1000,
                    700
                )

                return

            }
        }

        const size = getWorkingAreaSize()

        sizeAnimation(
            size.x,
            size.y,
            600,
            400
        )

    }

    let error = ""

    function saveLoadout(name: string) {

        if ( name.length <= 0 ) {
            error = "Name must have atleast 1 characters"
            console.log(error)
            return
        }

        error = ""

        SaveCurrentLoadout(name)

    }

    function deleteLoadout(name: string) {

        DeleteSavedLoadout(name)

    }

    function loadLoadout(name: string, isRandom: boolean) {

        error = "Loaded successfully!"
        LoadSavedLoadout(name)

    }

    function loadoutEqual(loadout1: LoadoutItem, loadout2: valorantapi.ValorantLocalLoadout ): boolean {

        if ( !loadout1 || !loadout2 ) {
            // Loadouts are not valid
            return false;
        }

        let GunsEqual = true

        for (let index = 0; index < loadout1.value.LoadoutData.Guns.length; index++) {
            const element1 = loadout1.value.LoadoutData.Guns[index];
            const element2 = loadout2.Guns[index];
            if (element1.ChromaID != element2.ChromaID) {
                GunsEqual = false
            }
        }

        let ExpressionsEqual = true

        for (let index = 0; index < loadout1.value.LoadoutData.ActiveExpressions.length; index++) {
            const element1 = loadout1.value.LoadoutData.ActiveExpressions[index];
            const element2 = loadout2.ActiveExpressions[index];
            if (element1.AssetID != element2.AssetID) {
                ExpressionsEqual = false
            }
        }

        return GunsEqual && ExpressionsEqual
         
    }
    
    type LoadoutItem = { key: string; value: main.SavedLoadout };
    var loadoutList: LoadoutItem[] = [];
    var CurrentLoadout: valorantapi.ValorantLocalLoadout

    onMount( () => {

        resize()

        EventsOn( "on_loadout_update", ( data: main.UpdateLoadoutObj ) => {

            CurrentLoadout = data.CurrentLoadout

            console.log("Current",CurrentLoadout)

            loadoutList = Object.entries(data.Loadouts).map( ([key, value]) => {
                return {
                    key: key,
                    value: value
                }
            })

            console.log(data)

        })

        GetLoadouts()

        return () => {

            EventsOff('on_loadout_update');

        }

    })

    let loadoutNameTextarea = '';
</script>

<main>

    <div class="container"
        bind:this={transferNode}>

        <div class="container-side">
            <bar>
                {#if error != ""}
                    <div class="container-text">{error}</div>
                {/if}
                {#if error == ""}
                    <div class="container-text">Loadouts</div>
                {/if}
                <textarea bind:value={loadoutNameTextarea} placeholder="Loadout Name" class="container-textarea"></textarea>
                <button on:click={ () => {saveLoadout(loadoutNameTextarea)}}>Save</button>
            </bar>
            
            <div class="loadout-list">
                {#each loadoutList as loadout (loadout.key) }
                    <div 
                        class="loadout-item {loadoutEqual( loadout, CurrentLoadout) ? "selected" : ""}" 
                        on:click={ () => { showLoadout(true, loadout) }}
                        on:keyup={ () => { showLoadout( true, loadout )}}
                    >
                        <img src="https://media.valorant-api.com/playercards/{loadout.value.LoadoutData.Identity.PlayerCardID}/displayicon.png" alt="Player Card" class="loadout-img"/>
                        <div class="card-holder">

                            <div class="vertical-devider"></div>

                            <card class="card">
                                <top>Loadout Name</top>
                                <div class="horizonal-devider"></div>
                                <bottom>{loadout.key}</bottom>
                            </card>
                            
                            <div class="vertical-devider"></div>

                        </div>
                    </div>
                {/each}
            </div>

        </div>

        <div class="container-side">
            
            {#if loadoutSelected != null}
            <bar>
                {#if error != ""}
                    <div class="container-text">{error}</div>
                {/if}
                {#if error == ""}
                    <div class="container-text">Loadouts</div>
                {/if}
                {#if loadoutSelected?.key != null}
                    <div class="container-title">{loadoutSelected.key}</div>
                {/if}
                <button on:click={ () => {loadLoadout(loadoutSelected?.key, false)}}>Load Loadout</button>
                <button on:click={ () => {
                    deleteLoadout(loadoutSelected.key)
                    showLoadout(false, null)
                }
                }>Delete</button>
                <button on:click={ () => {showLoadout(false, null)}}>Back</button>
            </bar>

                <div class="loadout-container">
                    
                    <div class="profile-card">
                        <img class="player-card" src="https://media.valorant-api.com/playercards/{loadoutSelected.value.LoadoutData.Identity.PlayerCardID}/displayicon.png" alt="Player card"/>
                        <div class="expression_container">
                            <div>Title</div>
                            <div class="player-title-container">
                                <div class="player-title">{loadoutSelected.value.NameLookup[loadoutSelected.value.LoadoutData.Identity.PlayerTitleID]}</div>
                            </div>
                        </div>
                        <div class="expression_container">
                            <div>Expressions</div>

                            <div class="expressions">

                                {#each loadoutSelected.value.LoadoutData.ActiveExpressions as item (item.AssetID)}

                                    {#if item.TypeID == "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475"}
                                        <img 
                                            class="expression" 
                                            src="https://media.valorant-api.com/sprays/{item.AssetID}/animation.png"
                                            {...{ onerror: `this.src='https://media.valorant-api.com/sprays/${item.AssetID}/fulltransparenticon.png'` }}
                                            alt="Spray"
                                        />
                                    {/if}
                                    {#if item.TypeID == "03a572de-4234-31ed-d344-ababa488f981"}
                                        <img 
                                            class="expression" 
                                            src="https://media.valorant-api.com/flex/{item.AssetID}/displayicon.png" 
                                            alt="Flex"
                                        />
                                    {/if}

                                {/each}
                                
                            </div>

                        </div>

                    </div>

                    <div class="skin_loadout">
                    
                        {#each loadoutSelected.value.LoadoutData.Guns as item (item.ID)}
                            
                            <div class="skin_loadout_item">

                                <div class="loadout_item_text">{loadoutSelected.value.NameLookup[item.ID]}</div>

                                <div class="image-container">

                                    <img src="https://media.valorant-api.com/weaponskinchromas/{item.ChromaID}/fullrender.png" alt="{item.ID}">
                                    
                                    {#if item.CharmID}
                                    <img class="loadout_item_buddy" style="left: {BuddyOffsets[item?.ID] ?? 0}px" src="https://media.valorant-api.com/buddies/{item.CharmID}/displayicon.png" alt="{item.CharmID}"/>
                                    {/if}
                                    
                                </div>

                            </div>

                        {/each}

                    </div>

                </div>

            {/if}

        </div>

    </div>

</main>

<style>

    .loadout-container {
        display: flex;
        flex-wrap: wrap;
        flex-direction: column;

        width: 100%;
        height: 100%;
    }

    .expression_container {
        display: row;
        justify-content: center;

        color: hsla(180, 67%, 99%, 0.7);
        font-weight: 700;
        font-family: 'DMSans', sans-serif;
    }

    .expressions {

        user-select: none;

        background-color: hsl(0, 0%, 10%);
        box-shadow: 0px 0px 0px 1px hsla(180, 67%, 99%, 0.5);
        height: 4rem;
        width: fit-content;
        padding: 0.25rem;
        border-radius: 1rem;

        display: flex;
        justify-content: center;
        align-items: center;
    }
    
    .expression {
        width: 4rem;
        height: 4rem;

        border-radius: 1rem;
    }

    .profile-card {
        display: flex;

        flex-direction: row;

        align-items: center;
        justify-content: center;

        gap: 2rem;
        
        width: 100%;
    }

    .selected {
        background-color: rgb(12, 34, 68) !important;
    }

    .player-card {
        user-select: none;
        width: auto;
        height: 7rem;

        border-radius: 1rem;
        box-shadow: 0px 0px 0px 2px hsla(180, 67%, 99%, 0.5);
    }

    .player-title-container {
        user-select: none;
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0px 0px 0px 1px hsla(180, 67%, 99%, 0.5);
        width: fit-content;
        padding: 0.4rem;
        border-radius: 1rem;

        margin-top: 1rem;

        box-sizing: border-box;

        line-height: 1rem;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .player-title {
        font-size: 0.8rem;

        font-weight: 300;
    }

    .skin_loadout {
        width: 100%;
        height: calc((78px + 0.6rem) * 6);

        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(230px, max-content));
        gap: 0.4rem;

        overflow-y: scroll;

        justify-content: center;
        padding-top: 2rem;
        padding-bottom: 1rem;

        box-sizing: border-box;
        scrollbar-width: none;
        -ms-overflow-style: none;
    }

    .skin_loadout::-webkit-scrollbar {
        display: none;
    }

    .skin_loadout_item {
        display: flex;
        flex-direction: column;
        align-items: center;

        user-select: none;

        gap: 0.5rem;

        width: auto;
        height: fit-content;

        border-radius: 0.25rem;
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0px 0px 0px 1px hsla(180, 67%, 99%, 0.5);
    }
    
    .image-container {
        position: relative;
        display: inline-block;
    }

    .skin_loadout_item img {
        height: 3rem;
        width: auto;
        object-fit: contain;
        margin-bottom: 0.25rem;
    }

    .loadout_item_buddy {

        position: absolute;
        animation: wobble 2s ease-in-out infinite;
        
        top: 13px;

        width: 2.5rem !important;
        height: auto !important;
        aspect-ratio: 1;

    }
    
    @keyframes wobble {
      0%   { transform: rotate(-5deg); }
      50%  { transform: rotate(5deg); }
      100% { transform: rotate(-5deg); }
    }

    .loadout_item_text {
        color: hsla(180, 67%, 99%, 0.7);
        font-weight: 700;
        font-family: 'DMSans', sans-serif;
        
        white-space: nowrap;
        text-overflow: ellipsis;
        overflow: hidden;
        width: 100%;
    }

    .container-side bar {
        display: flex;
        gap: 0.25rem;
        margin-top: 0.15rem;

        justify-content: end;
        align-items: center;

        width: 98%;
        height: 1.5rem;

        margin-left: 1%;
        margin-bottom: 0.5rem;

    }

    .container-text {

        all: unset;
        text-align: left;

        border: none;
        background-color: transparent;

        color: white;
        
        height: 1.5rem;

        padding-left: 1rem;
        padding-right: 1rem;

        flex: 1;
            
        font-family: 'DMSans', sans-serif;
        font-weight: 300;
        font-size: 0.8rem;

        text-align: center;
        line-height: 1.5rem;

        border-radius: 4px;
        
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);
    }

    .container-textarea {

        all: unset;
        text-align: center;

        border: none;
        background-color: transparent;

        color: white;
        
        width: fit-content;
        height: 1.5rem;

        border-radius: 4px;
        
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);
        
    }

    .container-title {
        border: none;
        background-color: transparent;

        color: white;
        
        width: fit-content;
        padding-left: 4rem;
        padding-right: 4rem;
        height: 1.5rem;

        border-radius: 4px;
        
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);
    }
    
    .container-side button {
        border: none;
        background-color: transparent;

        color: white;
        
        cursor: pointer;
        width: fit-content;
        height: 1.5rem;

        border-radius: 4px;

        transition: background-color 250ms cubic-bezier(0.87, 0, 0.13, 1);
        
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);
    } 

    .vertical-devider {
        align-self: stretch;
        width: 1px;
        height: auto;
        background-color: rgba(255, 255, 255, 0.25);
        margin-left: 0.125rem;
        margin-right: 0.125rem;
        flex-shrink: 0;
        
        align-self: stretch;
    }

    .horizonal-devider {
        align-self: stretch;
        width: auto;
        height: 1px;
        background-color: rgba(255, 255, 255, 0.25);
        margin-left: 0.125rem;
        margin-right: 0.125rem;
        flex-shrink: 0;
        
        align-self: stretch;
    }

    .card-holder {
        display: flex;

        justify-content: center;
        align-items: center;

        flex: 0.85;
    }

    .card {
        display: flex;

        margin-left: 1rem;
        margin-right: 1rem;

        align-items: center;

        flex-direction: column;
            
        font-family: 'DMSans', sans-serif;
    }

    .card top {
        
        font-weight: 700;
        font-size: 0.8rem;

        color: hsla(180, 67%, 99%, 0.5);

    }

    .card bottom {
        
        font-weight: 300;
        font-size: 1rem;

        align-self: center;

        color: hsla(180, 67%, 99%, 0.9);

    }

    .container {

        position: absolute;
        display: flex;

        width: 200%;
        height: 100%;
        
        transition: transform 300ms cubic-bezier(0.87, 0, 0.13, 1);

    }

    .container-side {
        flex: 1;
        width: 50%;
        height: 100%;
    }

    .loadout-list {

        user-select: none;
        
        display: flex;
        
        flex-direction: column;

        align-items: center;

        gap: 0.25rem;

        overflow-y: scroll;
        height: calc(360px - 2rem);
    }

    ::-webkit-scrollbar {
        display: none;
    }

    .loadout-item {
        display: flex;
        flex-direction: row;
        
        cursor: pointer;

        width: 95%;

        justify-content: start;
        align-items: center;

        gap: 1rem;

        padding-top: 0.5rem;
        padding-bottom: 0.5rem;

        background-color: hsl(0, 0%, 10%);
        border-radius: 1rem;

    }

    .loadout-img {
        margin-left: 1rem;

        width: auto;
        height: 2rem;
        
        border-radius: 0.5rem;
        box-shadow: 0px 0px 0px 2px hsla(180, 67%, 99%, 0.1);
    }

</style>