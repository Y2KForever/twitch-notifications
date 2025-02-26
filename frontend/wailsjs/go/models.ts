export namespace main {
	
	export class FileCheckResult {
	    exists: boolean;
	    content: string;
	    error: string;
	    errorType: string;
	
	    static createFrom(source: any = {}) {
	        return new FileCheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exists = source["exists"];
	        this.content = source["content"];
	        this.error = source["error"];
	        this.errorType = source["errorType"];
	    }
	}

}

