<fieldset>
	<label>
	<input type = "radio" name = "OS_Select" value = "Linux" 
	<?php if ($H_OSSel == "Linux") {echo "checked";}?>/>
	<span>Linux</span>
	</label>

	<label>
	<input type = "radio" name = "OS_Select" value = "Windows" 
	<?php if ($H_OSSel == "Windows") {echo "checked";}?>/>
	<span>Windows</span>
	</label>

	<label>
	<input type = "radio" name = "OS_Select" value = "Android" disabled 
	<?php if ($H_OSSel == "Android") {echo "checked";}?>/>
	<span>Android</span>
	</label>

	<label>
	<input type = "radio" name = "OS_Select" value = "temp" disabled
	<?php if ($H_OSSel == "temp") {echo "checked";}?>/>
	<span>Temp</span>
	</label>
</fieldset>	
