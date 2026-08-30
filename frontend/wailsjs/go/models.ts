export namespace project {
	
	export class Project {
	    id: string;
	    name: string;
	    path: string;
	    framework: string;
	    webRoot: string;
	    domain: string;
	    vhostEnabled: boolean;
	    sslEnabled: boolean;
	    targetPort?: number;
	    phpVersion?: string;
	    nodeVersion?: string;
	    createdAt?: string;
	    updatedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.framework = source["framework"];
	        this.webRoot = source["webRoot"];
	        this.domain = source["domain"];
	        this.vhostEnabled = source["vhostEnabled"];
	        this.sslEnabled = source["sslEnabled"];
	        this.targetPort = source["targetPort"];
	        this.phpVersion = source["phpVersion"];
	        this.nodeVersion = source["nodeVersion"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

export namespace service {
	
	export class Config {
	    displayName: string;
	    serviceName: string;
	    defaultPort: number;
	    category: string;
	    dependencies: string[];
	    healthCheckType: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.serviceName = source["serviceName"];
	        this.defaultPort = source["defaultPort"];
	        this.category = source["category"];
	        this.dependencies = source["dependencies"];
	        this.healthCheckType = source["healthCheckType"];
	    }
	}
	export class LogMessage {
	    timestamp: string;
	    level: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LogMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.level = source["level"];
	        this.message = source["message"];
	    }
	}
	export class RuntimeServiceInfo {
	    serviceName: string;
	    displayName: string;
	    activeVersion: string;
	    availableVersions: string[];
	    category: string;
	    isRunning: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeServiceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serviceName = source["serviceName"];
	        this.displayName = source["displayName"];
	        this.activeVersion = source["activeVersion"];
	        this.availableVersions = source["availableVersions"];
	        this.category = source["category"];
	        this.isRunning = source["isRunning"];
	    }
	}
	export class Status {
	    name: string;
	    isRunning: boolean;
	    message: string;
	    port?: number;
	    healthy?: boolean;
	    category?: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.isRunning = source["isRunning"];
	        this.message = source["message"];
	        this.port = source["port"];
	        this.healthy = source["healthy"];
	        this.category = source["category"];
	    }
	}

}

