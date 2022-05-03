using System;
using System.IO;

namespace VersionMakerNet
{
    class Program
    {
        static string VersionFilepath = "Version.txt";
        static string VersionLastFilepath = "VersionLast.txt";
        static string CommitHashFilepath = "CommitHash.txt";

        static ParameterParser Params = new ParameterParser();
        static FVersion Version = new FVersion();

        static void SaveCommitHash(string InHash)
        {
            File.WriteAllText(CommitHashFilepath, InHash);
        }

        static int CmdRaise(string InValue)
        {
            if (Params.Set.Length > 0)
            {
                Console.WriteLine("-Raise는 -Set과 동시에 사용할 수 없어 무시됩니다.");
                return 0;
            }

            string OldVersionString = Version.ToString();

            int VersionIndex = Int32.Parse(InValue);
            if (VersionIndex < 0)
            {
                Console.WriteLine("잘못된 매개변수입니다 : {0}", InValue);
                return -1;
            }

            Version.Raise(VersionIndex);

            string NewVersionString = Version.ToString();
            Console.WriteLine("버전이 변경됨 : {0} -> {1}", OldVersionString, NewVersionString);
            return 0;
        }

        static int CmdSet(string InValue)
        {
            string OldVersionString = Version.ToString();

            bool bSet = Version.SetByString(InValue);

            string NewVersionString = Version.ToString();

            if (bSet == false)
            {
                Console.WriteLine("버전 설정에 실패했습니다 : {0}", InValue);
                return -1;
            }

            Console.WriteLine("버전이 설정됨 : {0} -> {1}", OldVersionString, NewVersionString);

            return 0;
        }

        static int CmdCommitHash(string InValue)
        {
            if (InValue.Length <= 0)
            {
                Console.WriteLine("잘못된 매개변수입니다 : {0}", InValue);
                return -1;
            }

            string OldCommitHash = "";
            if (File.Exists(CommitHashFilepath))
                OldCommitHash = File.ReadAllText(CommitHashFilepath);
                
            if (OldCommitHash.Length <= 0)
            {
                SaveCommitHash(InValue);
                return 0;
            }

            Console.WriteLine("기존 커밋 : {0}", OldCommitHash);
            Console.WriteLine("현재 커밋 : {0}", InValue);

            int Index = -1;
            if (Params.CommitVer.Length > 0)
            {
                int CommitVer = Int32.Parse(Params.CommitVer);
                if (CommitVer >= 0 && CommitVer < FVersion.MAX_VERSION_COLUMN)
                    Index = CommitVer;
            }

            if (Index >= 0 && Index < FVersion.MAX_VERSION_COLUMN)
            {
                if (OldCommitHash.Equals(InValue))
                {
                    string VersionPosString = FVersion.GetVersionPosString(Index);
                    Console.WriteLine("커밋이 변경되지 않아 {0} 버전이 변경되지 않음.", VersionPosString);
                }
                else
                {
                    string OldVersionString = Version.ToString();
                    Version.Raise(Index);
                    string NewVersionString = Version.ToString();

                    SaveCommitHash(InValue);
                    Console.WriteLine("커밋이 변경되어 버전이 변경됨 : {0} -> {1}", OldVersionString, NewVersionString);

                }
            }
            else
            {
                SaveCommitHash(InValue);
                Console.WriteLine("버전을 올리지 않고 커밋 해쉬만 저장함.");
            }

            return 0;
        }

        static int CmdGenerate(string InValue)
        {
            if (InValue.Length <= 0)
                InValue = ".\\";

            bool bGenerated = Version.GenerateSourceFiles(InValue);
            if (!bGenerated)
            {
                Console.WriteLine("소스 파일을 생성할 수 없습니다 : {0}", InValue);
                return -1;
            }

            Console.WriteLine("{0} 에 .h, .cpp 생성 완료.", InValue);
            return 0;
        }

        // 제대로 하려면 INI 라입을 써야하지만, 변수 하나 수정하는 정도니 그냥 만듦.
        // 경고 : 해당 ini의 인코딩이 달라질 경우 이 기능은 제대로 작동하지 않을 수 있음.
        static int CmdGameIni(string InValue)
        {
            if (InValue.Length <= 0)
            {
                Console.WriteLine("잘못된 매개변수입니다 : {0}", InValue);
                return -1;
            }

            Console.WriteLine(InValue);

            // Load file
            if (!File.Exists(InValue))
            {
                Console.WriteLine("파일을 열 수 없습니다.");
                return -1;
            }

            string Buffer = "";
            if (File.Exists(InValue))
                Buffer = File.ReadAllText(InValue);

            string TargetVarName = "ProjectVersion";
            int FoundTargetVersionIndex = Buffer.IndexOf(TargetVarName);
            if (FoundTargetVersionIndex == -1)
            {
                Console.WriteLine("ini 파일 내에서 {0} 변수를 찾을 수 없습니다.", TargetVarName);
                return -1;
            }

            int FoundEqual = Buffer.IndexOf('=', FoundTargetVersionIndex);
            int FoundNewLine = Buffer.IndexOf('\n', FoundTargetVersionIndex);
            bool IsInvalidIndex = (FoundEqual != -1) && (FoundNewLine != -1);
            if (IsInvalidIndex || FoundNewLine < FoundEqual)
            {
                Console.WriteLine("ini 파일 내에서 {0} 변수를 찾을 수 없습니다.", TargetVarName);
                return -1;
            }

            int VersionValueIndex = FoundEqual + 1;

            int VersionLength = FoundEqual - VersionValueIndex;
            string IniVersionString = Buffer.Substring(VersionValueIndex, VersionLength);

            string VersionString = Version.ToString();
            string NewVersionStringFull = TargetVarName + "=" + VersionString + "\n";
            string IniVersionStringFull = Buffer.Substring(FoundTargetVersionIndex, FoundNewLine - FoundTargetVersionIndex);

            Buffer.Replace(IniVersionStringFull, NewVersionStringFull);

            File.WriteAllText(InValue, Buffer);

            Console.WriteLine("{0} 변경 완료 : {1} -> {2}", TargetVarName, IniVersionString, VersionString);
            return 0;
        }
        static int Main(string[] args)
        {
            // 1) parse argument
            Params.Parse(args);

            bool IsManualVersion = (Params.Set.Length > 0);

            if (Params.ChDir.Length > 0)
            {
                //JaeHwan Yi, 맥에서 끝에 슬래쉬가 있으면 Directory.Exists가 실패함.
                if (Params.ChDir.EndsWith('/') || Params.ChDir.EndsWith('\\'))
                {
                    Params.ChDir = Params.ChDir.Remove(Params.ChDir.Length - 1);

                    if (!Directory.Exists(Params.ChDir))
                    {
                        Directory.CreateDirectory(Params.ChDir);
                    }
                }
                else
                {
                    Directory.CreateDirectory(Params.ChDir);
                }

                Directory.SetCurrentDirectory(Params.ChDir);
            }

            if (Params.DataDir.Length > 0)
            {
                //JaeHwan Yi, 맥에서 끝에 슬래쉬가 있으면 Directory.Exists가 실패함.
                if (Params.DataDir.EndsWith('/') || Params.DataDir.EndsWith('\\'))
                {
                    Params.DataDir = Params.DataDir.Remove(Params.DataDir.Length - 1);
                    if (!Directory.Exists(Params.DataDir))
                    {
                        Directory.CreateDirectory(Params.DataDir);
                    }
                }
                else
                {
                    Directory.CreateDirectory(Params.DataDir);
                }

                VersionFilepath = Path.Combine(Params.DataDir, VersionFilepath);
                VersionLastFilepath = Path.Combine(Params.DataDir, VersionLastFilepath);
                CommitHashFilepath = Path.Combine(Params.DataDir, CommitHashFilepath);
            }

            // 2) load version from file
            bool bNeedsToSave = false;
            if (!IsManualVersion) //자동 업데이트
            {
                if (0 < Params.Redis.Length) // 레디스 있으면
                {
                    bool bLoaded = Version.LoadFromRedis(Params.Redis);
                    if (!bLoaded)
                    {
                        Console.WriteLine("레디스 서버로부터 버전을 불러오는데 실패했습니다.");
                        return 1;
                    }
                    bNeedsToSave = true;
                }
                else
                {
                    bool bLoaded = Version.LoadFromFile(VersionFilepath); //local version file 
                    if (!bLoaded)
                        bNeedsToSave = true;
                }
            }

            FVersion VersionLast = Version.Clone();
            string LastVersionString = VersionLast.ToString();

            // 3) print version and out if it is just single argument
            if (args.Length == 1)
            {
                Console.WriteLine(LastVersionString);
                return 1;
            }

            // 4) process all registered argument(functions)
            int NumOfProcessed = 0;

            if (Params.Raise.Length > 0)
            {
                CmdRaise(Params.Raise);
                ++NumOfProcessed;
            }
            if (Params.Set.Length > 0)
            {
                CmdSet(Params.Set);
                ++NumOfProcessed;
            }
            if (Params.CommitHash.Length > 0)
            {
                CmdCommitHash(Params.CommitHash);
                ++NumOfProcessed;
            }
            if (Params.Generate.Length > 0)
            {
                CmdGenerate(Params.Generate);
                ++NumOfProcessed;
            }
            if (Params.GameIni.Length > 0)
            {
                CmdGameIni(Params.GameIni);
                ++NumOfProcessed;
            }

            // 5) print out the final version
            string NewVersionString = Version.ToString();

            bool bPrintLastVersion = Boolean.Parse(Params.Last);

            if (NumOfProcessed == 0 || Boolean.Parse(Params.IsQuiet) || bPrintLastVersion)
            {
                if (bPrintLastVersion)
                {
                    Console.WriteLine(LastVersionString);
                }
                else
                {
                    Console.WriteLine(NewVersionString);
                }
            }
            else
            {
                Console.WriteLine("");
                Console.WriteLine("기존 버전 : {0}", LastVersionString);
                Console.WriteLine("현재 버전 : {0}", NewVersionString);

                // 6) save changed version 
                if (bNeedsToSave == false && !VersionLast.Equals(Version))
                    bNeedsToSave = true;

                if (bNeedsToSave)
                {
                    VersionLast.SaveToFile(VersionLastFilepath);
                    Version.SaveToFile(VersionFilepath);
                    bNeedsToSave = false;

                    Console.WriteLine("저장됨.");
                }

                Console.WriteLine("");
                Console.WriteLine("작업 완료.");
            }

            return 0;
        }
    }
}