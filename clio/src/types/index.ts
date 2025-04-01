export interface Folder {
  id: number;
  name: string;
  count: number;
  resourceUrl: string;
  lastSynced: string;
  createdAt: string;
  updatedAt: string;
}

export interface Artist {
  id: number;
  name: string;
  resourceUrl: string;
}

export interface ReleaseArtist {
  releaseId: number;
  artistId: number;
  joinRelation: string;
  anv: string;
  tracks: string;
  role: string;
  artist?: Artist;
}

export interface Label {
  id: number;
  name: string;
  resourceUrl: string;
  entityType: string;
}

export interface ReleaseLabel {
  releaseId: number;
  labelId: number;
  catNo: string;
  label?: Label;
}

export interface Format {
  id: number;
  releaseId: number;
  name: string;
  qty: number;
  descriptions: string[];
}

export interface Genre {
  id: number;
  name: string;
}

export interface Style {
  id: number;
  name: string;
}

export interface ReleaseNote {
  releaseId: number;
  fieldId: number;
  value: string;
}

export interface Release {
  id: number;
  instanceId: number;
  folderId: number;
  rating: number;
  title: string;
  year: number | null;
  resourceUrl: string;
  thumb: string;
  coverImage: string;
  createdAt: string;
  updatedAt: string;
  lastSynced: string;
  playDuration: number;
  playDurationEstimated: boolean;

  artists: ReleaseArtist[];
  labels: ReleaseLabel[];
  formats: Format[];
  genres: Genre[];
  styles: Style[];
  notes: ReleaseNote[];
  playHistory: PlayHistory[];
  cleaningHistory: CleaningHistory[];
}

export interface Stylus {
  id: number;
  name: string;
  manufacturer?: string;
  expectedLifespan?: number;
  purchaseDate?: string;
  active: boolean;
  primary: boolean;
  createdAt: string;
  updatedAt: string;
  owned: boolean;
  baseModel: boolean;
}

export interface PlayHistory {
  id: number;
  releaseId: number;
  stylusId?: number;
  playedAt: string;
  createdAt: string;
  updatedAt: string;
  release?: Release;
  stylus?: Stylus;
  owned: boolean;
  notes?: string;
}
export interface CleaningHistory {
  id: number;
  releaseId: number;
  cleanedAt: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
  release?: Release;
}

export interface EditItem {
  id: number;
  type: "play" | "cleaning";
  date: Date;
  notes?: string;
  stylus?: string;
  stylusId?: number;
  releaseId: number;
}
