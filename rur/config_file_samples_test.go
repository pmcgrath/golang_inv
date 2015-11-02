package main

var configContentEmpty string = ``

var configContentNoContent string = `<?xml version="1.0" encoding="utf-8"?><configuration></configuration>`

var configContentSample string = `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <configSections>
    <section name="nlog" type="NLog.Config.ConfigSectionHandler, NLog" />
  </configSections>
  <nlog>
    <extensions>
      <add assembly="NLog.Targets.Gelf" />
    </extensions>
    <targets>
      <target name="UdpOutlet" type="NLogViewer" address="udp://localhost:7071" />
      <target name="Gelf" type="Gelf" facility="ServiceABC" gelfserver="log.company.com" port="12201" maxchunksize="8154" graylogversion="0.9.6" />
    </targets>
    <rules>
      <logger name="*" minLevel="Trace" writeTo="UdpOutlet" />
      <logger name="*" minLevel="Trace" appendTo="Gelf" />
    </rules>
  </nlog>
  <appSettings>
    <add key="IncludeJsonTypeInfo" value="true" />
    <add key="Metrics.GlobalContextName" value="ServiceABC" />
    <add key="SecurityServiceUrl" value="http://dev.security.com/" />
  </appSettings>
  <connectionStrings>
    <add name="AAADatabase" connectionString="Data Source=db1; Initial Catalog=AAA; MultiSubnetFailover=True; Integrated Security=SSPI;" providerName="System.Data.SqlClient" />
    <add name="BBBDatabase" connectionString="Data Source=db1; Initial Catalog=BBB; Integrated Security=SSPI;" providerName="System.Data.SqlClient" />
    <add name="CCDatabase" connectionString="Data Source=db2; Initial Catalog=CCC; MultiSubnetFailover=True; Integrated Security=SSPI;" providerName="System.Data.SqlClient" />
    <add name="EventStore" connectionString="CONNECTTo=tcp://admin:ASuperDupperStrongPassword@127.0.0.1:1113" />,
  </connectionStrings>
  <rabbitServers>
    <add key="mb1" value="Host=Rabbit1" />
  </rabbitServers>
  <system.web>
    <compilation debug="true" targetFramework="4.5" />
    <httpRuntime targetFramework="4.5" />
  </system.web>
  <runtime>
    <assemblyBinding xmlns="urn:schemas-microsoft-com:asm.v1">
      <dependentAssembly>
        <assemblyIdentity name="Castle.Core" publicKeyToken="407dd0808d44fbdc" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-3.3.0.0" newVersion="3.3.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="Newtonsoft.Json" publicKeyToken="30ad4fe6b2a6aeed" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-6.0.0.0" newVersion="6.0.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="Microsoft.Owin" publicKeyToken="31bf3856ad364e35" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-2.1.0.0" newVersion="2.1.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="NLog" publicKeyToken="5120e14c03d0593c" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-3.1.0.0" newVersion="3.1.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="Castle.Services.Logging.NLogIntegration" publicKeyToken="407dd0808d44fbdc" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-3.3.0.0" newVersion="3.3.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="Castle.Windsor" publicKeyToken="407dd0808d44fbdc" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-3.3.0.0" newVersion="3.3.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="System.Web.Http" publicKeyToken="31bf3856ad364e35" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-5.2.0.0" newVersion="5.2.0.0" />
      </dependentAssembly>
      <dependentAssembly>
        <assemblyIdentity name="System.Net.Http.Formatting" publicKeyToken="31bf3856ad364e35" culture="neutral" />
        <bindingRedirect oldVersion="0.0.0.0-5.2.0.0" newVersion="5.2.0.0" />
      </dependentAssembly>
    </assemblyBinding>
  </runtime>
  <system.webServer>
    <handlers>
      <remove name="ExtensionlessUrlHandler-Integrated-4.0" />
      <remove name="OPTIONSVerbHandler" />
      <remove name="TRACEVerbHandler" />
      <add name="ExtensionlessUrlHandler-Integrated-4.0" path="*." verb="*" type="System.Web.Handlers.TransferRequestHandler" preCondition="integratedMode,runtimeVersionv4.0" />
    </handlers>
  </system.webServer>
</configuration>
`

var configContentTransformation string = `<?xml version="1.0" encoding="utf-8"?>
<?xml version="1.0"?>
<configuration xmlns:xdt="http://schemas.microsoft.com/XML-Document-Transform">
  <nlog>
    <targets>
      <target name="Gelf" gelfserver="log.company.com" xdt:Locator="Match(name)" xdt:Transform="SetAttributes" />
    </targets>
    <rules xdt:Transform="Replace">
      <logger name="*" minLevel="Warn" appendTo="GelfOld" />
      <logger name="*" minLevel="Debug" appendTo="Gelf" />
    </rules>
  </nlog>
  <appSettings>
    <add key="IsLive" value="True" xdt:Transform="Insert" />
    <add key="IncludeJsonTypeInfo" value="false" xdt:Locator="Match(key)" xdt:Transform="SetAttributes" />
    <add key="SecurityServiceUrl" value="https://security.services.local/api2/security/" xdt:Locator="Match(key)" xdt:Transform="SetAttributes" />
  </appSettings>
  <connectionStrings xdt:Transform="Replace">
    <add name="AAADatabase" connectionString="Server=tcp:db2; Database=AAA; MultiSubnetFailover=True; Integrated Security=SSPI;" providerName="System.Data.SqlClient" />
    <add name="BBBDatabase" connectionString="Data Source=db2;Initial Catalog=BBB; Integrated Security=SSPI;" providerName="System.Data.SqlClient" />
    <add name="CCCDatabase" connectionString="Data Source=cdb2; Initial Catalog=CCC; Integrated Security=SSPI;" providerName="System.Data.SqlClient" />
    <add name="metrics" connectionString="host=metrics01;port=8086;user=root;password=PASSWORD;database=metrics"/>
  </connectionStrings>
</configuration>
`
