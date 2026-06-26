export namespace history {
	
	export class Message {
	    id: string;
	    project_id: string;
	    phase: string;
	    role: string;
	    content: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.phase = source["phase"];
	        this.role = source["role"];
	        this.content = source["content"];
	        this.created_at = source["created_at"];
	    }
	}
	export class Project {
	    id: string;
	    name: string;
	    tagline: string;
	    path: string;
	    provider: string;
	    model: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.tagline = source["tagline"];
	        this.path = source["path"];
	        this.provider = source["provider"];
	        this.model = source["model"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}

}

export namespace providers {
	
	export class ChatResponse {
	    content: string;
	    model: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.model = source["model"];
	    }
	}
	export class Message {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class Model {
	    id: string;
	    display_name: string;
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.display_name = source["display_name"];
	    }
	}

}

export namespace settings {
	
	export class ProviderConfig {
	    endpoint: string;
	    default_model: string;
	
	    static createFrom(source: any = {}) {
	        return new ProviderConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpoint = source["endpoint"];
	        this.default_model = source["default_model"];
	    }
	}
	export class Config {
	    providers: Record<string, ProviderConfig>;
	    default_provider: string;
	    appearance: string;
	    default_license: string;
	    default_project_dir: string;
	    font: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.providers = this.convertValues(source["providers"], ProviderConfig, true);
	        this.default_provider = source["default_provider"];
	        this.appearance = source["appearance"];
	        this.default_license = source["default_license"];
	        this.default_project_dir = source["default_project_dir"];
	        this.font = source["font"];
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

export namespace templates {
	
	export class ProjectFiles {
	    ProjectMD: string;
	    AgentsMD: string;
	    PlanMD: string;
	    ReadmeMD: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectFiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ProjectMD = source["ProjectMD"];
	        this.AgentsMD = source["AgentsMD"];
	        this.PlanMD = source["PlanMD"];
	        this.ReadmeMD = source["ReadmeMD"];
	    }
	}

}

