<script type="ts">

    import { onMount } from 'svelte';
    import { EventsOn, EventsOff } from '../../wailsjs/runtime';
    import { valorantapi, main } from '../../wailsjs/go/models';
    import { SaveCurrentLoadout, DeleteSavedLoadout, LoadSavedLoadout, GetLoadouts, AddWeaponsToBeRandomized, GetRandomLoadout, LoadRandomLoadout } from '../../wailsjs/go/main/App';
    
    import { 
        sizeAnimation,
        getWorkingAreaSize
    } from './../utils/animations';
    import { element } from 'svelte/internal';

    let transferNode: HTMLElement;
    let isLoadoutShown: boolean = false;

    let loadoutSelected: LoadoutItem
    
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
        }

    }
    
    function resize() {

        const size = getWorkingAreaSize()

        sizeAnimation(
            size.x,
            size.y,
            1000,
            700
        )

    }

    let error = ""

    function loadRandom(isRandom: boolean) {

        error = "Loaded successfully!"
        LoadRandomLoadout(isRandom)

    }
    
    type LoadoutItem = { key: string; value: main.SavedLoadout };
    var randomLoadout: main.SavedLoadout
    var CurrentLoadout: valorantapi.ValorantLocalLoadout
    let isRandomEnabled: boolean = false;
    let randomizedItems: Record<string, boolean>;

    function setItemToRandomize( ID: string ) {

        if (randomizedItems[ID] == null) {
            randomizedItems[ID] = true
        } else {
            randomizedItems[ID] = !randomizedItems[ID]
        }

        AddWeaponsToBeRandomized( randomizedItems)

    }

    onMount( () => {

        resize()

        EventsOn( "on_random_update", ( data: main.UpdateRandomObj ) => {

            CurrentLoadout = data.CurrentLoadout
            randomizedItems = data.RandomWeaponsSelected

            console.log("Current",CurrentLoadout)

            data.RandomLoadout.LoadoutData.Identity.PreferredLevelBorderID = CurrentLoadout.Identity.PreferredLevelBorderID

            data.RandomLoadout.LoadoutData.Identity.PlayerTitleID = "random"
            data.RandomLoadout.NameLookup[data.RandomLoadout.LoadoutData.Identity.PlayerTitleID] = "Random Title"

            data.RandomLoadout.LoadoutData.ActiveExpressions = CurrentLoadout.ActiveExpressions

            randomLoadout = data.RandomLoadout
            isRandomEnabled = data.IsRandomSelected

            console.log(data)

        })

        GetRandomLoadout()

        return () => {

            EventsOff('on_random_update');

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
                    <div class="container-text">Skin Randomizer</div>
                {/if}
                {#if loadoutSelected?.key != null}
                    <div class="container-title">{loadoutSelected.key}</div>
                {/if}
                <button class="{ isRandomEnabled ? "selected" : "" }" on:click={ () => {loadRandom(!isRandomEnabled )}}>{ isRandomEnabled ? "Disable" : "Enable" } Random</button>
            </bar>
            
            {#if randomLoadout != null}

            <div class="loadout-container">
                
                <div class="profile-card">
                    <img class="player-card" src="https://media.valorant-api.com/playercards/{randomLoadout?.LoadoutData?.Identity?.PlayerCardID}/displayicon.png" alt="Player card"/>
                    <div class="expression_container">
                        <div>Title</div>
                    <div class="player-title-container">
                        <div class="player-title">{randomLoadout?.NameLookup[randomLoadout?.LoadoutData?.Identity.PlayerTitleID]}</div>
                    </div>
                    </div>
                    <div class="expression_container">
                        <div>Expressions</div>

                        <div class="expressions">
                            {#each randomLoadout?.LoadoutData?.ActiveExpressions as item (item.AssetID)}
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
                
                    {#each CurrentLoadout?.Guns as item, index (item.ID)}
                        
                        <div 
                            class="skin_loadout_item { randomizedItems[item.ID] ? "selected" : "" }"
                            on:click={ () => { setItemToRandomize(item.ID) }}
                            on:keyup={ () => { setItemToRandomize(item.ID) }}
                        >   

                            <div class="loadout_item_text">{randomLoadout.NameLookup[item.ID]}</div>
                            <img src="https://media.valorant-api.com/weaponskinchromas/{item.ChromaID}/fullrender.png" alt="{item.ID}"/>

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

    .skin_loadout_item img {
        height: 3rem;
        width: 225px;
        object-fit: contain;
        margin-bottom: 0.25rem;
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
        margin-top: 0.25rem;

        justify-content: end;
        align-items: center;

        width: 98%;
        height: 1.5rem;

        margin-left: 1%;
        margin-bottom: 1rem;

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

    .container {

        position: absolute;

        width: 100%;
        height: 100%;

        display: flex;
        
        transition: transform 300ms cubic-bezier(0.87, 0, 0.13, 1);

    }

    .container-side {
        width: 100%;
        height: 100%;
    }

    .loadout-list {

        user-select: none;
        
        display: flex;
        
        flex-direction: column;

        justify-content: center;
        align-items: center;

        gap: 0.25rem;

        height: fit-content;
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