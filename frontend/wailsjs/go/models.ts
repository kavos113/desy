export namespace domain {
	
	export class LecturePlan {
	    Count: number;
	    Plan: string;
	    Assignment: string;
	
	    static createFrom(source: any = {}) {
	        return new LecturePlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Count = source["Count"];
	        this.Plan = source["Plan"];
	        this.Assignment = source["Assignment"];
	    }
	}
	export class Teacher {
	    ID: number;
	    Name: string;
	    Url: string;
	
	    static createFrom(source: any = {}) {
	        return new Teacher(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Url = source["Url"];
	    }
	}
	export class Room {
	    ID: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Room(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	    }
	}
	export class TimeTable {
	    LectureID: number;
	    Semester: string;
	    Room: Room;
	    DayOfWeek: string;
	    Period: number;
	
	    static createFrom(source: any = {}) {
	        return new TimeTable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LectureID = source["LectureID"];
	        this.Semester = source["Semester"];
	        this.Room = this.convertValues(source["Room"], Room);
	        this.DayOfWeek = source["DayOfWeek"];
	        this.Period = source["Period"];
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
	export class Lecture {
	    ID: number;
	    University: string;
	    Title: string;
	    EnglishTitle: string;
	    Department: string;
	    LectureType: string;
	    Code: string;
	    Level: number;
	    Credit: number;
	    Year: number;
	    OpenTerm: string;
	    Language: string;
	    Url: string;
	    Abstract: string;
	    Goal: string;
	    Experience: string;
	    Flow: string;
	    OutOfClassWork: string;
	    Textbook: string;
	    ReferenceBook: string;
	    Assessment: string;
	    Prerequisite: string;
	    Contact: string;
	    OfficeHours: string;
	    Note: string;
	    // Go type: time
	    UpdatedAt: any;
	    Timetables: TimeTable[];
	    Teachers: Teacher[];
	    LecturePlans: LecturePlan[];
	    Keywords: string[];
	    RelatedCourseCodes: string[];
	    RelatedCourses: number[];
	
	    static createFrom(source: any = {}) {
	        return new Lecture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.University = source["University"];
	        this.Title = source["Title"];
	        this.EnglishTitle = source["EnglishTitle"];
	        this.Department = source["Department"];
	        this.LectureType = source["LectureType"];
	        this.Code = source["Code"];
	        this.Level = source["Level"];
	        this.Credit = source["Credit"];
	        this.Year = source["Year"];
	        this.OpenTerm = source["OpenTerm"];
	        this.Language = source["Language"];
	        this.Url = source["Url"];
	        this.Abstract = source["Abstract"];
	        this.Goal = source["Goal"];
	        this.Experience = source["Experience"];
	        this.Flow = source["Flow"];
	        this.OutOfClassWork = source["OutOfClassWork"];
	        this.Textbook = source["Textbook"];
	        this.ReferenceBook = source["ReferenceBook"];
	        this.Assessment = source["Assessment"];
	        this.Prerequisite = source["Prerequisite"];
	        this.Contact = source["Contact"];
	        this.OfficeHours = source["OfficeHours"];
	        this.Note = source["Note"];
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.Timetables = this.convertValues(source["Timetables"], TimeTable);
	        this.Teachers = this.convertValues(source["Teachers"], Teacher);
	        this.LecturePlans = this.convertValues(source["LecturePlans"], LecturePlan);
	        this.Keywords = source["Keywords"];
	        this.RelatedCourseCodes = source["RelatedCourseCodes"];
	        this.RelatedCourses = source["RelatedCourses"];
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
	
	export class LectureSummary {
	    ID: number;
	    University: string;
	    Title: string;
	    Department: string;
	    Code: string;
	    Level: number;
	    Credit: number;
	    Year: number;
	    Timetables: TimeTable[];
	    Teachers: Teacher[];
	
	    static createFrom(source: any = {}) {
	        return new LectureSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.University = source["University"];
	        this.Title = source["Title"];
	        this.Department = source["Department"];
	        this.Code = source["Code"];
	        this.Level = source["Level"];
	        this.Credit = source["Credit"];
	        this.Year = source["Year"];
	        this.Timetables = this.convertValues(source["Timetables"], TimeTable);
	        this.Teachers = this.convertValues(source["Teachers"], Teacher);
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
	
	export class SearchQuery {
	    Title: string;
	    Keywords: string[];
	    Departments: string[];
	    Year: number;
	    TeacherName: string;
	    TimeTables: TimeTable[];
	    Levels: number[];
	
	    static createFrom(source: any = {}) {
	        return new SearchQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Keywords = source["Keywords"];
	        this.Departments = source["Departments"];
	        this.Year = source["Year"];
	        this.TeacherName = source["TeacherName"];
	        this.TimeTables = this.convertValues(source["TimeTables"], TimeTable);
	        this.Levels = source["Levels"];
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

