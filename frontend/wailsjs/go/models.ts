export namespace client {
	
	export class CreateRequest {
	    name: string;
	    weight: number;
	    height: number;
	    target: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.weight = source["weight"];
	        this.height = source["height"];
	        this.target = source["target"];
	    }
	}
	export class Response {
	    id: string;
	    name: string;
	    weight: number;
	    height: number;
	    target: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.weight = source["weight"];
	        this.height = source["height"];
	        this.target = source["target"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class UpdateRequest {
	    id: string;
	    name: string;
	    weight: number;
	    height: number;
	    target: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.weight = source["weight"];
	        this.height = source["height"];
	        this.target = source["target"];
	    }
	}

}

export namespace day_variant {
	
	export class CreateRequest {
	    planId: string;
	    name: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.planId = source["planId"];
	        this.name = source["name"];
	        this.order = source["order"];
	    }
	}
	export class Response {
	    id: string;
	    planId: string;
	    name: string;
	    order: number;
	    meals?: meal.Response[];
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.planId = source["planId"];
	        this.name = source["name"];
	        this.order = source["order"];
	        this.meals = this.convertValues(source["meals"], meal.Response);
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
	export class UpdateRequest {
	    id: string;
	    name: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.order = source["order"];
	    }
	}

}

export namespace food {
	
	export class CreateRequest {
	    mealId: string;
	    baseFood: string;
	    baseQuantity: number;
	    baseUnit: string;
	    alternativeFood?: string;
	    alternativeQuantity?: number;
	    alternativeUnit?: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mealId = source["mealId"];
	        this.baseFood = source["baseFood"];
	        this.baseQuantity = source["baseQuantity"];
	        this.baseUnit = source["baseUnit"];
	        this.alternativeFood = source["alternativeFood"];
	        this.alternativeQuantity = source["alternativeQuantity"];
	        this.alternativeUnit = source["alternativeUnit"];
	        this.order = source["order"];
	    }
	}
	export class Response {
	    id: string;
	    mealId: string;
	    baseFood: string;
	    baseQuantity: number;
	    baseUnit: string;
	    alternativeFood?: string;
	    alternativeQuantity?: number;
	    alternativeUnit?: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.mealId = source["mealId"];
	        this.baseFood = source["baseFood"];
	        this.baseQuantity = source["baseQuantity"];
	        this.baseUnit = source["baseUnit"];
	        this.alternativeFood = source["alternativeFood"];
	        this.alternativeQuantity = source["alternativeQuantity"];
	        this.alternativeUnit = source["alternativeUnit"];
	        this.order = source["order"];
	    }
	}
	export class UpdateRequest {
	    id: string;
	    baseFood: string;
	    baseQuantity: number;
	    baseUnit: string;
	    alternativeFood?: string;
	    alternativeQuantity?: number;
	    alternativeUnit?: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.baseFood = source["baseFood"];
	        this.baseQuantity = source["baseQuantity"];
	        this.baseUnit = source["baseUnit"];
	        this.alternativeFood = source["alternativeFood"];
	        this.alternativeQuantity = source["alternativeQuantity"];
	        this.alternativeUnit = source["alternativeUnit"];
	        this.order = source["order"];
	    }
	}

}

export namespace meal {
	
	export class CreateRequest {
	    dayVariantId: string;
	    name: string;
	    note: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dayVariantId = source["dayVariantId"];
	        this.name = source["name"];
	        this.note = source["note"];
	        this.order = source["order"];
	    }
	}
	export class Response {
	    id: string;
	    dayVariantId: string;
	    name: string;
	    note: string;
	    order: number;
	    foods?: food.Response[];
	    supplements?: supplement.Response[];
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.dayVariantId = source["dayVariantId"];
	        this.name = source["name"];
	        this.note = source["note"];
	        this.order = source["order"];
	        this.foods = this.convertValues(source["foods"], food.Response);
	        this.supplements = this.convertValues(source["supplements"], supplement.Response);
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
	export class UpdateRequest {
	    id: string;
	    name: string;
	    note: string;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.note = source["note"];
	        this.order = source["order"];
	    }
	}

}

export namespace pdf_export {
	
	export class Response {
	    fileName: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileName = source["fileName"];
	        this.path = source["path"];
	    }
	}

}

export namespace plan {
	
	export class CreateRequest {
	    clientId: string;
	    title: string;
	    generalRecommendations: string;
	    footer: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.title = source["title"];
	        this.generalRecommendations = source["generalRecommendations"];
	        this.footer = source["footer"];
	    }
	}
	export class FullResponse {
	    id: string;
	    clientId: string;
	    title: string;
	    generalRecommendations: string;
	    footer: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    dayVariants?: day_variant.Response[];
	
	    static createFrom(source: any = {}) {
	        return new FullResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.title = source["title"];
	        this.generalRecommendations = source["generalRecommendations"];
	        this.footer = source["footer"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.dayVariants = this.convertValues(source["dayVariants"], day_variant.Response);
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
	export class Response {
	    id: string;
	    clientId: string;
	    title: string;
	    generalRecommendations: string;
	    footer: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.title = source["title"];
	        this.generalRecommendations = source["generalRecommendations"];
	        this.footer = source["footer"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class UpdateRequest {
	    id: string;
	    title: string;
	    generalRecommendations: string;
	    footer: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.generalRecommendations = source["generalRecommendations"];
	        this.footer = source["footer"];
	    }
	}

}

export namespace settings {
	
	export class Response {
	    exportBaseDirectory: string;
	    colorPageBackground: string;
	    colorCardBackground: string;
	    colorCardBorder: string;
	    colorPrimary: string;
	    colorPrimaryLight: string;
	    colorOnPrimary: string;
	    colorText: string;
	    colorMuted: string;
	    colorTableHeader: string;
	    colorDivider: string;
	    uiTheme: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exportBaseDirectory = source["exportBaseDirectory"];
	        this.colorPageBackground = source["colorPageBackground"];
	        this.colorCardBackground = source["colorCardBackground"];
	        this.colorCardBorder = source["colorCardBorder"];
	        this.colorPrimary = source["colorPrimary"];
	        this.colorPrimaryLight = source["colorPrimaryLight"];
	        this.colorOnPrimary = source["colorOnPrimary"];
	        this.colorText = source["colorText"];
	        this.colorMuted = source["colorMuted"];
	        this.colorTableHeader = source["colorTableHeader"];
	        this.colorDivider = source["colorDivider"];
	        this.uiTheme = source["uiTheme"];
	        this.language = source["language"];
	    }
	}
	export class UpdateRequest {
	    exportBaseDirectory?: string;
	    colorPageBackground?: string;
	    colorCardBackground?: string;
	    colorCardBorder?: string;
	    colorPrimary?: string;
	    colorPrimaryLight?: string;
	    colorOnPrimary?: string;
	    colorText?: string;
	    colorMuted?: string;
	    colorTableHeader?: string;
	    colorDivider?: string;
	    uiTheme?: string;
	    language?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exportBaseDirectory = source["exportBaseDirectory"];
	        this.colorPageBackground = source["colorPageBackground"];
	        this.colorCardBackground = source["colorCardBackground"];
	        this.colorCardBorder = source["colorCardBorder"];
	        this.colorPrimary = source["colorPrimary"];
	        this.colorPrimaryLight = source["colorPrimaryLight"];
	        this.colorOnPrimary = source["colorOnPrimary"];
	        this.colorText = source["colorText"];
	        this.colorMuted = source["colorMuted"];
	        this.colorTableHeader = source["colorTableHeader"];
	        this.colorDivider = source["colorDivider"];
	        this.uiTheme = source["uiTheme"];
	        this.language = source["language"];
	    }
	}

}

export namespace supplement {
	
	export class CreateRequest {
	    mealId: string;
	    supplement: string;
	    quantity: number;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mealId = source["mealId"];
	        this.supplement = source["supplement"];
	        this.quantity = source["quantity"];
	        this.order = source["order"];
	    }
	}
	export class Response {
	    id: string;
	    mealId: string;
	    supplement: string;
	    quantity: number;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.mealId = source["mealId"];
	        this.supplement = source["supplement"];
	        this.quantity = source["quantity"];
	        this.order = source["order"];
	    }
	}
	export class UpdateRequest {
	    id: string;
	    supplement: string;
	    quantity: number;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.supplement = source["supplement"];
	        this.quantity = source["quantity"];
	        this.order = source["order"];
	    }
	}

}

