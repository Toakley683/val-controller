<script lang="ts">

    import { onMount } from 'svelte';
    import { EventsOn, EventsOff } from '../../wailsjs/runtime';
    import { valorantapi } from '../../wailsjs/go/models';
    import { ExitCoreGame, ExitPregame, GetMatch, SelectRandomAgent } from '../../wailsjs/go/main/App';

    import { 
        sizeAnimation,
        getWorkingAreaSize
    } from './../utils/animations';

    let MatchData: valorantapi.MatchData
    let transferNode: HTMLElement;

    let maxPlayerPerSize = 0
    let hasBeenUpdated = false

    let isInventoryShown: boolean = false;
    let selectedPlayer: valorantapi.valorantMatchTeamPlayer;

    const BuddyOffsets = {
        "63e6c2b6-4a8e-869c-3d4c-e38355226584": 70,
        "55d8a0f4-4274-ca67-fe2c-06ab45efdf58": 50,
        "9c82e19d-4575-0200-1a81-3eacf00cf872": 25,
        "ae3de142-4d85-2547-dd26-4e90bed35cf7": 35,
        "ee8e8d15-496b-07ac-e5f6-8fae5d4c7b1a": 45,
        "ec845bf4-4f79-ddda-a3da-0db3774b2794": 30,
        "910be174-449b-c412-ab22-d0873436b21b": 45,
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

    let TeamLookup = {}
    
    function objectToArray<T>(obj: Record<string, T>): T[] {
        return Object.values(obj);
    }

    function showInventory( isShown: boolean, player: valorantapi.valorantMatchTeamPlayer ) {

        if (transferNode == null) {
            return
        }

        isInventoryShown = isShown

        if (isShown ) {
            transferNode.style.transform = "translate(-50%, 0px)";
        } else {
            transferNode.style.transform = "translate(0%, 0px)";
        }

        if (player != null) {
            selectedPlayer = player
        }

        resize()

    }

    onMount( () => {

        showInventory( false, null )

        EventsOn( "update_match", ( data: valorantapi.MatchData ) => {

            hasBeenUpdated = true

            console.log(data)

            maxPlayerPerSize = 0

            MatchData = data

            TeamLookup = {}

            data.AllyTeam?.Players?.forEach(element => {

                if (TeamLookup[element.LastMatchPartyID] == null){
                    TeamLookup[element.LastMatchPartyID] = 0
                }

                TeamLookup[element.LastMatchPartyID] = TeamLookup[element.LastMatchPartyID] + 1
                
            });

            data.EnemyTeam?.Players?.forEach(element => {

                if (TeamLookup[element.LastMatchPartyID] == null){
                    TeamLookup[element.LastMatchPartyID] = 0
                }

                TeamLookup[element.LastMatchPartyID] = TeamLookup[element.LastMatchPartyID] + 1
                
            });

            console.log("Team Lookup:", TeamLookup)

            if (selectedPlayer == null) {
                
                selectedPlayer = data.AllyTeam.Players[0]

            }
            
            console.log("Resize")

            resize()

        })

        GetMatch()

        return () => {

            EventsOff('update_match');

        }

    })

    function resize() {

        if (transferNode != null ) {
            if ( isInventoryShown ) {

                const size = getWorkingAreaSize()

                sizeAnimation(
                    size.x,
                    size.y,
                    1000,
                    600
                )

                return

            }
        }

        const size = getWorkingAreaSize()

        if (MatchData == null) {

            sizeAnimation(
                size.x,
                size.y,
                550,
                41 + 32
            )
            return
        }

        if (MatchData.AllyTeam.Players != null) {
            maxPlayerPerSize = MatchData.AllyTeam.Players.length
        }

        if (MatchData.EnemyTeam.Players != null) {

            if (maxPlayerPerSize < MatchData.EnemyTeam.Players.length) {

                maxPlayerPerSize = MatchData.EnemyTeam.Players.length

            }
        }

        let GapAdd = 0

        if (maxPlayerPerSize > 1) {
            GapAdd = 4
        }

        if (maxPlayerPerSize > 0) {

            if (MatchData.IsPregame) {

            
                sizeAnimation(
                    size.x,
                    size.y,
                    600,
                    40 + 20 + 16 + ( ( 47 + GapAdd ) * maxPlayerPerSize )
                )
                
            } else {

                sizeAnimation(
                    size.x,
                    size.y,
                    600,
                    40 + 10 + ( ( 47 + GapAdd ) * maxPlayerPerSize )
                )

            }

        } else {

            sizeAnimation(
                size.x,
                size.y,
                600,
                41 + 32
            )

        }

    }

    function djb2(str){
    var hash = 5381;
    for (var i = 0; i < str.length; i++) {
        hash = ((hash << 5) + hash) + str.charCodeAt(i); /* hash * 33 + c */
    }
    return hash;
    }

    function hashStringToColor(str) {
        var hash = djb2(str);
        var r = (hash & 0xFF0000) >> 16;
        var g = (hash & 0x00FF00) >> 8;
        var b = hash & 0x0000FF;
        return "#" + ("0" + r.toString(16)).substr(-2) + ("0" + g.toString(16)).substr(-2) + ("0" + b.toString(16)).substr(-2);
    }

</script>

<main>

    <div class="transfer-node"
        bind:this={transferNode}>

        <div class="transfer-node-side">

            {#if maxPlayerPerSize <= 0 && hasBeenUpdated}

                <div class="match-reminder">Enter Match To Start</div>

            {/if}

            {#if MatchData?.AllyTeam?.Players?.length > 0}
            
                {#if MatchData?.IsPregame }
                    <bar>
                        <button on:click={ () => { ExitPregame() } }>Exit Pregame</button>
                        <button on:click={ () => { SelectRandomAgent() } }>Random Agent</button>
                    </bar>
                {/if}

            {/if}

            <div class="team-sides">

                <div class="ally-side">
                
                    {#if MatchData?.AllyTeam?.Players?.length > 0}

                        {#each MatchData.AllyTeam.Players as player (player.Subject)}
                            
                            <div class="agent-row" 
                                on:click={() => {showInventory( true, player )}}
                                on:keyup={() => {showInventory( true, player )}}
                            >
                                                    
                                {#if TeamLookup[player.LastMatchPartyID] > 1}

                                    <!--/* Only if in party */-->
                                    <div class="party-indicator" style="background-color: {hashStringToColor(player.LastMatchPartyID)}"></div>

                                {/if}
                                
                                {#if TeamLookup[player.LastMatchPartyID] <= 1}

                                    <!--/* Only if in party */-->
                                    <div class="party-indicator"></div>

                                {/if}

                                <boundary class="boundary">

                                    <img src="{player.CharacterDisplayIcon}" alt="{player.CharacterName}">
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>
                                            Username
                                            <div class="horizonal-devider"></div>
                                        </top>
                                        <bottom>
                                            {player.PlayerIdentity.GameName}
                                            {#if !player.PlayerIdentity.Incognito }
                                                <top>#{player.PlayerIdentity.TagLine}</top>
                                            {/if}
                                        </bottom>
                                    </card>

                                </boundary>

                                <boundary class="boundary-reverse">
                                    
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>Peak</top>
                                        <div class="horizonal-devider"></div>
                                        <img src={player.PeakRankDisplayIcon} alt={player.PeakRank}>
                                    </card>
                                    
                                    <div class="vertical-devider"></div>
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>Rank</top>
                                        <div class="horizonal-devider"></div>
                                        <img src={player.CurrentRankDisplayIcon} alt={player.CurrentRank}>
                                    </card>
                                    
                                    <div class="vertical-devider"></div>
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>Level</top>
                                        <div class="horizonal-devider"></div>
                                        
                                        {#if !player.PlayerIdentity.HideAccountLevel}
                                            <bottom>{player.PlayerIdentity.AccountLevel}</bottom>
                                        {/if}
                                        
                                        {#if player.PlayerIdentity.HideAccountLevel }
                                            <bottom>?</bottom>
                                        {/if}

                                    </card>
                                    
                                    <div class="vertical-devider"></div>

                                </boundary>

                            </div>

                        {/each}

                    {/if}
                    
                </div>
                
                {#if MatchData?.EnemyTeam?.Players?.length > 0}

                    <div class="enemy-side">

                        {#each MatchData.EnemyTeam.Players as player (player.Subject)}
                            
                            <div class="agent-row"
                                on:click={() => {showInventory( true, player )}}
                                on:keyup={() => {showInventory( true, player )}}
                            >
                                
                                {#if TeamLookup[player.LastMatchPartyID] > 1}

                                    <!--/* Only if in party */-->
                                    <div class="party-indicator" style="background-color: {hashStringToColor(player.LastMatchPartyID)}"></div>

                                {/if}
                                
                                {#if TeamLookup[player.LastMatchPartyID] <= 1}

                                    <!--/* Only if in party */-->
                                    <div class="party-indicator"></div>

                                {/if}

                                <boundary class="boundary">

                                    <img src="{player.CharacterDisplayIcon}" alt="{player.CharacterName}">
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>
                                            Username
                                            <div class="horizonal-devider"></div>
                                        </top>
                                        <bottom>
                                            {player.PlayerIdentity.GameName}
                                            {#if !player.PlayerIdentity.Incognito }
                                                <top>#{player.PlayerIdentity.TagLine}</top>
                                            {/if}
                                        </bottom>
                                    </card>

                                </boundary>

                                <boundary class="boundary-reverse">
                                    
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>Peak</top>
                                        <div class="horizonal-devider"></div>
                                        <img src={player.PeakRankDisplayIcon} alt={player.PeakRank}>
                                    </card>
                                    
                                    <div class="vertical-devider"></div>
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>Rank</top>
                                        <div class="horizonal-devider"></div>
                                        <img src={player.CurrentRankDisplayIcon} alt={player.CurrentRank}>
                                    </card>
                                    
                                    <div class="vertical-devider"></div>
                                    <div class="vertical-devider"></div>
                                    
                                    <card>
                                        <top>Level</top>
                                        <div class="horizonal-devider"></div>
                                        
                                        {#if !player.PlayerIdentity.HideAccountLevel}
                                            <bottom>{player.PlayerIdentity.AccountLevel}</bottom>
                                        {/if}
                                        
                                        {#if player.PlayerIdentity.HideAccountLevel }
                                            <bottom>?</bottom>
                                        {/if}
                                    </card>
                                    
                                    <div class="vertical-devider"></div>

                                </boundary>

                            </div>


                        {/each}

                    </div>

                {/if}

            </div>

            </div>

        <div 
            class="transfer-node-side inventory-side">
            
            <bar>
                <button on:click={ () => {showInventory(false, null)}}>Back</button>
            </bar>
            
            {#if selectedPlayer != null}

            <profile_card>
                <img class="player_icon" src="{selectedPlayer.CharacterDisplayIcon}" alt="{selectedPlayer.CharacterName}"/>

                <div class="vertical-devider"></div>

                <card>
                    <top>
                        Username
                        <div class="horizonal-devider"></div>
                    </top>
                    <bottom>
                        {selectedPlayer.PlayerIdentity.GameName}
                        {#if !selectedPlayer.PlayerIdentity.Incognito }
                            <top>#{selectedPlayer.PlayerIdentity.TagLine}</top>
                        {/if}
                    </bottom>
                </card>

                <div class="vertical-devider"></div>
                
                {#if selectedPlayer.MatchesAgo > 0 }

                    <card>
                        <top>
                            Last Seen
                            <div class="horizonal-devider"></div>
                        </top>
                        <bottom>
                            {#if selectedPlayer.MatchesAgo > 1 }
                                {selectedPlayer.MatchesAgo} Matches ago
                            {/if}
                            {#if selectedPlayer.MatchesAgo == 1 }
                                {selectedPlayer.MatchesAgo} Match ago
                            {/if}
                        </bottom>
                    </card>

                    <div class="vertical-devider"></div>

                {/if}

                <card>
                    <top>Level</top>
                    <div class="horizonal-devider"></div>
                    
                    {#if !selectedPlayer.PlayerIdentity.HideAccountLevel}
                        <bottom>{selectedPlayer.PlayerIdentity.AccountLevel}</bottom>
                    {/if}
                    
                    {#if selectedPlayer.PlayerIdentity.HideAccountLevel }
                        <bottom>?</bottom>
                    {/if}
                </card>

                <div class="vertical-devider"></div>

                <card>
                    <top>Rank</top>
                    <div class="horizonal-devider"></div>
                    <img src={selectedPlayer.CurrentRankDisplayIcon} alt={selectedPlayer.CurrentRank}>
                </card>
                
                <card>
                    <top>Peak</top>
                    <div class="horizonal-devider"></div>
                    <img src={selectedPlayer.PeakRankDisplayIcon} alt={selectedPlayer.PeakRank}>
                </card>
                
                <div class="vertical-devider"></div>
                

                <div class="expression_container">
                        <top>Expressions</top>
                        <div class="horizonal-devider"></div>
                        
                        <div class="expressions">

                        {#each selectedPlayer.Items.ActiveExpressions as item (item.AssetID)}

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

            </profile_card>

            <!-- Match History -->

            <!-- Skin Loadout -->

            <div class="skin_loadout">
                
                {#each selectedPlayer.Items.Guns as item, index (item.ID)}
                
                <div class="skin_loadout_item">

                    <div class="loadout_item_text">{item.SkinName}</div>

                    <div class="image-container">

                        <img src="https://media.valorant-api.com/weaponskinchromas/{item.ChromaID}/fullrender.png" alt="{item.SkinName}">
                        
                        {#if item.CharmID}
                        <img class="loadout_item_buddy" style="left: {BuddyOffsets[item?.ID] ?? 0}px" src="https://media.valorant-api.com/buddies/{item.CharmID}/displayicon.png" alt="{item.CharmID}"/>
                        {/if}
                        
                    </div>

                </div>

                {/each}

            </div>

            {/if}
        </div>

    </div>

</main>

<style>

    .transfer-node {
        position: absolute;

        width: 200%;
        height: fit-content;

        display: flex;
        justify-content: start;
        
        
        transition: transform 300ms cubic-bezier(0.87, 0, 0.13, 1);
        
    }

    .transfer-node-side {
        margin: 0;
        flex: 1;
    }

    .skin_loadout {
        margin-left: 2.5%;
        width: 95%;
        height: 100%;

        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(230px, max-content));
        gap: 0.4rem;

        justify-content: center;
        padding-top: 1rem;
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
        height: 2.5rem;
        width: auto;
        object-fit: contain;
        margin-bottom: 0.25rem;
    }

    .loadout_item_buddy {

        position: absolute;
        animation: wobble 2s ease-in-out infinite;
        
        top: 10px;

        width: 2rem !important;
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

    .inventory-side {
        justify-content: start;
        
        display: flex;
        flex-direction: column;
    }

    .transfer-node bar {
        
        display: flex;
        gap: 0.25rem;
        margin-top: 0.25rem;

        justify-content: end;

        width: 98%;
        height: 1.5rem;

        margin-left: 1%;

    }

    button {
        border: none;
        background-color: transparent;

        color: white;
        
        cursor: pointer;
        width: fit-content;
        height: 1.5rem;

        border-radius: 4px;
        
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);
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
        height: 2rem;
        width: fit-content;
        padding: 0.25rem;
        border-radius: 1rem;

        display: flex;
        justify-content: center;
        align-items: center;
    }
    
    .expression {
        width: 2rem !important;
        height: 2rem !important;

        border-radius: 1rem;
    }

    .inventory-side profile_card {
        display: flex;

        gap: 10px;
        align-items: center;
        justify-content: center;

        user-select: none;

        width: 100%;

        margin-top: 1rem;
        height: 4rem;

        font-weight: 700;
        font-size: 1rem;
    }

    .inventory-side profile_card top {
        color: hsla(180, 67%, 99%, 0.5);
    }

    .inventory-side profile_card bottom {
        color: hsla(180, 67%, 99%, 0.9);
    }
    
    .inventory-side profile_card img {
        width: 1.5rem;
        height: 1.5rem;
    }

    .player_icon {

        height: 100% !important;
        width: auto !important;

        background-color: hsl(0, 0%, 10%);
        border-radius: 0.5rem;
        box-shadow: 0px 0px 0px 2px hsla(180, 67%, 99%, 0.1);
    }

    .inventory-side button {
        border: none;
        background-color: transparent;

        color: white;
        
        cursor: pointer;
        width: 4rem;
        height: 1.5rem;

        border-radius: 4px;
        
        background-color: hsl(0, 0%, 10%);
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);
    }

    .team-sides {

        display: flex;

        justify-content: center;
        align-items: center;

        gap: 15px;

        padding: 5px;

        width: 100%;
        height: 100%;

        box-sizing: border-box;

    }

    .team-sides > * {
        height: 100%;
        flex: 1;

        gap: 5px;

        display: flex;
        flex-wrap:  nowrap;
        flex-direction: column;
        justify-content: start;

    }

    .team-sides bar {
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

    .ally-side {

        /*background-color: rgb(57, 126, 190);*/

    }

    .enemy-side {

        /*background-color: rgb(190, 57, 57);*/

    }

    .agent-row {

        user-select: none;

        display: flex;
        flex-direction: row;
        gap: 0.25rem;
        padding: 0.425rem 0.55rem;

        cursor: pointer;
        overflow: hidden;

        justify-content: center;
        align-items: center;

        border-radius: 0.875rem;
        background-color: hsl(0, 0%, 10%);
        /*border: 1px solid hsla(180, 67%, 99%, 0.1);;*/
        box-shadow: 0 2px 0.5rem rgba(0, 0, 0, 0.3);

    }

    .agent-row img {

        height: 2rem;
        width: 2rem;
        
        border-radius: 0.2rem;
        /*border: 2px solid hsla(180, 67%, 99%, 0.1);*/
        box-shadow: 0px 0px 0px 1.5px hsla(180, 67%, 99%, 0.1);

        transform: scale(94%);

    }

    .party-indicator {
        align-self: flex-start;
        justify-self: left;
        transform: scaleY(200%);
        margin-left: -0.55rem;
        width: 4px;
        height: 100%;
        background-color: rgba(255, 255, 255, 0.25);
        flex-shrink: 0;
        
        align-self: stretch;
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

    .agent-row boundary card {
        display: flex;

        align-items: center;
        justify-content: center;

        flex: 1;

        flex-direction: column;
            
        font-family: 'DMSans', sans-serif;
    }

    .agent-row boundary card top {
        
        font-weight: 700;
        font-size: 0.6rem;

        color: hsla(180, 67%, 99%, 0.5);

    }

    .agent-row boundary card bottom {
        
        font-weight: 300;
        font-size: 0.7rem;

        align-self: center;

        color: hsla(180, 67%, 99%, 0.9);

    }

    .agent-row boundary card img {
        height: 1.3rem;
        width: 1.3rem;
    }

    .boundary {

        display: flex;

        flex: 1;

    }

    .boundary-reverse {
        display: flex;

        flex-direction: row-reverse;

    }

    .match-reminder {
        width: 100%;
        height: 16px;
        margin: 0;
        margin-top: 8px;
        
        font-family: 'DMSans', sans-serif;
        font-weight: 100;
        font-size: 1rem;
    }

</style>