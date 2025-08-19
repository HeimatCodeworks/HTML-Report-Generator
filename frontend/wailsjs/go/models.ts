export namespace main {
	
	export class Company {
	    teamId: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Company(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.teamId = source["teamId"];
	        this.name = source["name"];
	    }
	}
	export class ReportInfo {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new ReportInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

