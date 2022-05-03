using System;
using System.Collections.Generic;
using System.Reflection;
using System.Text;

namespace VersionMakerNet
{
    public class ParameterParser
    {
        // Function Arguments
        public string Raise = "";           // -raise=index 처럼 사용하며, 해당 인덱스의 버전을 1만큼 올립니다. 인덱스 값이 없을 경우 가장 뒷 자리의 버전을 올립니다.
        public string Set = "";      // -set=0.3.1.5 처럼 사용하며, 현재 버전을 설정합니다.
        public string CommitHash = "";  // 현재 커밋 해쉬를 저장합니다. 마지막으로 저장된 해쉬와 다를 경우, CommitVer 값에 해당하는 버전이 1만큼 올립니다.
        public string Generate = "";         // 현재 버전 정보로 .h 및 .cpp 파일을 생성합니다. -generate=PATH 로 생성할 경로를 지정할 수 있습니다.
        public string GameIni = "";         // Unreal Engine 4의 DefaultGame.ini의 경로를 입력합니다. ini 파일 내 ProjectVersion= 변수에 현재 버전 정보를 입력합니다.

        // Options
        public string IsQuiet = "false";           // 오류를 포함하여 작업하는 동안 출력을 하지 않습니다.
        public string ChDir = "";           // -chdir=path 처럼 사용하며, 작업 디렉토리를 지정합니다. 기본값은 exe 파일이 위치한 곳입니다.
        public string DataDir = "";         // -datadir=path 처럼 사용하며, 데이터 디렉토리를 지정합니다. 기본값은 exe 파일이 위치한 곳입니다.
        public string CommitVer = "";       // -commitver=index 처럼 사용하며, CommitHash가 변경됐다면 index에 해당하는 버전을 1만큼 올립니다. 기본값은 -1이며 버전을 올리지 않습니다.
        public string Last = "false";            // -last 처럼 사용하며, 현재 버전이 아닌 마지막 버전을 출력합니다.
        public string Redis = "";      // 레디스 서버를 통해 버전 수정

        public string Parse(string[] args)
        {
            try
            {
                Dictionary<string, System.Reflection.FieldInfo> FieldDictionary = new Dictionary<string, FieldInfo>();

                FieldInfo[] Fields = GetType().GetFields();
                foreach (FieldInfo field in Fields)
                {
                    // 내 클래스의 변수만 파싱함
                    if (field.DeclaringType.Name == MethodBase.GetCurrentMethod().DeclaringType.Name)
                        FieldDictionary.Add(field.Name.ToLower(), field);
                }

                string NotUsedArgList = "";

                Console.WriteLine("Parsing input arguments");
                for (int i = 0; i < args.Length; ++i)
                {
                    string[] SplipttedArg = args[i].Split('=');

                    if (SplipttedArg.Length == 0)
                        continue;

                    if (SplipttedArg.Length == 2)
                    {
                        string FieldName = SplipttedArg[0];
                        if (FieldName[0] == '-')
                            FieldName = FieldName.Substring(1);

                        FieldInfo info;
                        if (!FieldDictionary.TryGetValue(FieldName.ToLower(), out info))
                        {
                            NotUsedArgList += " " + args[i];
                            continue;
                        }

                        FieldDictionary.Remove(FieldName.ToLower());
                        info.SetValue(this, SplipttedArg[1]);
                        SetEnvironmentVariable(info.Name, SplipttedArg[1]);
                    }
                    else if (SplipttedArg.Length == 1)
                    {
                        string FieldName = SplipttedArg[0];
                        if (FieldName[0] == '-')
                            FieldName = FieldName.Substring(1);

                        FieldInfo info;
                        if (!FieldDictionary.TryGetValue(FieldName.ToLower(), out info))
                        {
                            NotUsedArgList += " " + args[i];
                            continue;
                        }

                        FieldDictionary.Remove(FieldName.ToLower());
                        info.SetValue(this, "true");
                        SetEnvironmentVariable(info.Name, "true");
                    }
                    else
                    {
                        Console.WriteLine("Invalid Args : ", args[i]);
                    }
                }

                // Argument로 입력되지 않는 내용들이 환경변수에 등록되어있다면, 그것을 사용하도록 함.
                Console.WriteLine("Setting parameterParser's variables form EnvironmentVariables");
                foreach (var FieldPair in FieldDictionary)
                {
                    FieldInfo info = FieldPair.Value;
                    string result = GetEnvironmentVariable(info.Name);
                    if (result.Length > 0)
                        info.SetValue(this, result);
                }

                return NotUsedArgList;
            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
            }

            return "";
        }

        public string GetEnvironmentVariable(string InName, string InDefault = "")
        {
            string Variable = Environment.GetEnvironmentVariable(InName);
            if (Variable == null)
                return InDefault;
            return Variable;
        }

        public void SetEnvironmentVariable(string InName, string InValue)
        {
            Environment.SetEnvironmentVariable(InName, InValue);
            Console.WriteLine("EnvironmentVariable Registered (" + InName + ":" + InValue + ")");
        }
    }
}