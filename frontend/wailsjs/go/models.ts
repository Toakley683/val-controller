export namespace main {
	
	export class SavedLoadout {
	    LoadoutData: valorantapi.ValorantLocalLoadout;
	    NameLookup: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new SavedLoadout(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LoadoutData = this.convertValues(source["LoadoutData"], valorantapi.ValorantLocalLoadout);
	        this.NameLookup = source["NameLookup"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateLoadoutObj {
	    Loadouts: Record<string, SavedLoadout>;
	    CurrentLoadout: valorantapi.ValorantLocalLoadout;
	
	    static createFrom(source: any = {}) {
	        return new UpdateLoadoutObj(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Loadouts = this.convertValues(source["Loadouts"], SavedLoadout, true);
	        this.CurrentLoadout = this.convertValues(source["CurrentLoadout"], valorantapi.ValorantLocalLoadout);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace valorantapi {
	
	export class valorantItemLoadout {
	    DisplayIcon: string;
	    DisplayName: string;
	
	    static createFrom(source: any = {}) {
	        return new valorantItemLoadout(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DisplayIcon = source["DisplayIcon"];
	        this.DisplayName = source["DisplayName"];
	    }
	}
	export class PlayerIdentity {
	    Subject: string;
	    PlayerCardID: string;
	    PlayerTitleID: string;
	    AccountLevel: number;
	    PreferredLevelBorderID: string;
	    Incognito: boolean;
	    HideAccountLevel: boolean;
	    GameName: string;
	    TagLine: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerIdentity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Subject = source["Subject"];
	        this.PlayerCardID = source["PlayerCardID"];
	        this.PlayerTitleID = source["PlayerTitleID"];
	        this.AccountLevel = source["AccountLevel"];
	        this.PreferredLevelBorderID = source["PreferredLevelBorderID"];
	        this.Incognito = source["Incognito"];
	        this.HideAccountLevel = source["HideAccountLevel"];
	        this.GameName = source["GameName"];
	        this.TagLine = source["TagLine"];
	    }
	}
	export class valorantMatchTeamPlayer {
	    Subject: string;
	    CharacterID: string;
	    CharacterName: string;
	    CharacterDisplayIcon: string;
	    CharacterSelectionState: string;
	    PeakRank: string;
	    PeakRankDisplayIcon: string;
	    CurrentRank: string;
	    CurrentRankDisplayIcon: string;
	    LastMatchPartyID: string;
	    MatchesAgo: number;
	    PregamePlayerState: string;
	    CompetitiveTier: number;
	    PlayerIdentity: PlayerIdentity;
	    IsCaptain: boolean;
	    PlatformType: string;
	    Items: Record<string, valorantItemLoadout>;
	
	    static createFrom(source: any = {}) {
	        return new valorantMatchTeamPlayer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Subject = source["Subject"];
	        this.CharacterID = source["CharacterID"];
	        this.CharacterName = source["CharacterName"];
	        this.CharacterDisplayIcon = source["CharacterDisplayIcon"];
	        this.CharacterSelectionState = source["CharacterSelectionState"];
	        this.PeakRank = source["PeakRank"];
	        this.PeakRankDisplayIcon = source["PeakRankDisplayIcon"];
	        this.CurrentRank = source["CurrentRank"];
	        this.CurrentRankDisplayIcon = source["CurrentRankDisplayIcon"];
	        this.LastMatchPartyID = source["LastMatchPartyID"];
	        this.MatchesAgo = source["MatchesAgo"];
	        this.PregamePlayerState = source["PregamePlayerState"];
	        this.CompetitiveTier = source["CompetitiveTier"];
	        this.PlayerIdentity = this.convertValues(source["PlayerIdentity"], PlayerIdentity);
	        this.IsCaptain = source["IsCaptain"];
	        this.PlatformType = source["PlatformType"];
	        this.Items = this.convertValues(source["Items"], valorantItemLoadout, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ValorantMatchTeam {
	    TeamID: string;
	    TeamNumber: number;
	    Players: valorantMatchTeamPlayer[];
	
	    static createFrom(source: any = {}) {
	        return new ValorantMatchTeam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TeamID = source["TeamID"];
	        this.TeamNumber = source["TeamNumber"];
	        this.Players = this.convertValues(source["Players"], valorantMatchTeamPlayer);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MatchData {
	    MatchID: string;
	    IsPregame: boolean;
	    AllyTeam: ValorantMatchTeam;
	    EnemyTeam: ValorantMatchTeam;
	    MapID: string;
	
	    static createFrom(source: any = {}) {
	        return new MatchData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MatchID = source["MatchID"];
	        this.IsPregame = source["IsPregame"];
	        this.AllyTeam = this.convertValues(source["AllyTeam"], ValorantMatchTeam);
	        this.EnemyTeam = this.convertValues(source["EnemyTeam"], ValorantMatchTeam);
	        this.MapID = source["MapID"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ValorantLocalExpression {
	    TypeID: string;
	    AssetID: string;
	
	    static createFrom(source: any = {}) {
	        return new ValorantLocalExpression(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TypeID = source["TypeID"];
	        this.AssetID = source["AssetID"];
	    }
	}
	export class ValorantLocalLoadoutIdentity {
	    PlayerCardID: string;
	    PlayerTitleID: string;
	    AccountLevel: number;
	    PreferredLevelBorderID: string;
	    Incognito: boolean;
	    HideAccountLevel: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ValorantLocalLoadoutIdentity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PlayerCardID = source["PlayerCardID"];
	        this.PlayerTitleID = source["PlayerTitleID"];
	        this.AccountLevel = source["AccountLevel"];
	        this.PreferredLevelBorderID = source["PreferredLevelBorderID"];
	        this.Incognito = source["Incognito"];
	        this.HideAccountLevel = source["HideAccountLevel"];
	    }
	}
	export class ValorantLocalLoadoutGuns {
	    ID: string;
	    SkinID: string;
	    SkinLevelID: string;
	    ChromaID: string;
	    CharmInstanceID?: string;
	    CharmID?: string;
	    CharmLevelID?: string;
	    Attachments: any[];
	
	    static createFrom(source: any = {}) {
	        return new ValorantLocalLoadoutGuns(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SkinID = source["SkinID"];
	        this.SkinLevelID = source["SkinLevelID"];
	        this.ChromaID = source["ChromaID"];
	        this.CharmInstanceID = source["CharmInstanceID"];
	        this.CharmID = source["CharmID"];
	        this.CharmLevelID = source["CharmLevelID"];
	        this.Attachments = source["Attachments"];
	    }
	}
	export class ValorantLocalLoadout {
	    Subject: string;
	    Version: number;
	    Guns: ValorantLocalLoadoutGuns[];
	    ActiveExpressions: ValorantLocalExpression[];
	    Identity: ValorantLocalLoadoutIdentity;
	
	    static createFrom(source: any = {}) {
	        return new ValorantLocalLoadout(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Subject = source["Subject"];
	        this.Version = source["Version"];
	        this.Guns = this.convertValues(source["Guns"], ValorantLocalLoadoutGuns);
	        this.ActiveExpressions = this.convertValues(source["ActiveExpressions"], ValorantLocalExpression);
	        this.Identity = this.convertValues(source["Identity"], ValorantLocalLoadoutIdentity);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	

}

