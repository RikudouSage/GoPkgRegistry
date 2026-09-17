export enum Vcs {
  Git = "git",
  Mercurial = "hg",
  Subversion = "svn",
  Fossil = "fossil",
}

export interface Package<TNullType = null> {
  id: number;

  import_path: string;
  vcs: Vcs;
  repository_url: string;

  source_url: string | TNullType;
  source_dir_url: string | TNullType;
  source_file_url: string | TNullType;
}

export const emptyPackage = <TNullType>(nullValue: TNullType): Package<TNullType> => ({
  id: 0,
  import_path: '',
  vcs: Vcs.Git,
  repository_url: '',
  source_url: nullValue,
  source_dir_url: nullValue,
  source_file_url: nullValue,
});
