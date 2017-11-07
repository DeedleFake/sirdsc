package sirdsc

import (
	"image"
	"image/color"
)

// A RandImage deterministically generates random color values for
// each (x, y) coordinate, using itself as a seed. In other words,
// given two RandImages that are equal to each other, the color at
// the same (x, y) in each are also equal.
//
// A RandImage has infinite size.
type RandImage uint64

func (img RandImage) ColorModel() color.Model { // nolint
	return color.RGBAModel
}

func (img RandImage) Bounds() image.Rectangle { // nolint
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img RandImage) At(x, y int) color.Color { // nolint
	base := uint64(int(img)^x) ^ rand[int(uint64(y)%uint64(len(rand)))]

	r := uint8(base ^ rand[int(base%uint64(len(rand)))])
	base ^= rand[int(img%RandImage(len(rand)))]

	g := uint8(base ^ rand[int(base%uint64(len(rand)))])
	base ^= rand[int(base%uint64(len(rand)))]

	b := uint8(base ^ rand[int(base%uint64(len(rand)))])

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

var rand = []uint64{
	5577006791947779410, 8674665223082153551, 15352856648520921629, 13260572831089785859, 3916589616287113937,
	6334824724549167320, 9828766684487745566, 10667007354186551956, 894385949183117216, 11998794077335055257,
	4751997750760398084, 7504504064263669287, 11199607447739267382, 3510942875414458836, 12156940908066221323,
	4324745483838182873, 11833901312327420776, 11926759511765359899, 6263450610539110790, 11239168150708129139,
	1874068156324778273, 3328451335138149956, 14486903973548550719, 7955079406183515637, 11926873763676642186,
	2740103009342231109, 6941261091797652072, 1905388747193831650, 17204678798284737396, 15649472107743074779,
	4831389563158288344, 261049867304784443, 10683692646452562431, 5600924393587988459, 18218388313430417611,
	9956202364908137547, 5486140987150761883, 9768663798983814715, 6382800227808658932, 2781055864473387780,
	10821471013040158923, 4990765271833742716, 14242321332569825828, 11792151447964398879, 13126262220165910460,
	14117161486975057715, 2338498362660772719, 2601737961087659062, 7273596521315663110, 3337066551442961397,
	17344948852394588913, 11963748953446345529, 8249030965139585917, 898860202204764712, 9010467728050264449,
	9908585559158765387, 11273630029763932141, 15505210698284655633, 2227583514184312746, 12096659438561119542,
	8603989663476771718, 6842348953158377901, 7388428680384065704, 6735196588112087610, 1687184559264975024,
	13174268766980400525, 17496662575514578077, 6296367092202729479, 18252401681137062077, 8505906760983331750,
	837825985403119657, 13771804148684671731, 8549944162621642512, 8807817071862113702, 12432680895096110463,
	15595235597337683065, 6556961545928831643, 5199948958991797301, 15213854965919594827, 5089134323978233018,
	16194613440650274502, 12947799971452915849, 10428415896243638596, 18317291550776694829, 17490665426807838719,
	2970700287221458280, 6651414131918424343, 5944830206637008055, 788787457839692041, 15399114114227588261,
	14967026985784794439, 3409814636252858217, 11407674492757219439, 4937104021912138218, 10950412492527322440,
	2202916659517317514, 5793183108815074904, 1169089424364679180, 11818186001859264308, 3784560248718450071,
	13234731587024579193, 14988856041258832631, 14297581759627478249, 5751776211841778805, 6725505124774569258,
	16883695360970880573, 9228111148518271676, 16012406593094538891, 12912571090385939658, 4592022834646721379,
	14794086776323742620, 13001433806883825921, 9497041302863216379, 9240932364278507230, 9249594463326629931,
	8446960703956728189, 14663632165210175172, 14382856709244076395, 16744157289148322445, 17321601139639006207,
	13451757574255826437, 5535550569387508244, 9465625292531964560, 17024802514613298337, 2303013289404122822,
	5919415281453547599, 2282476590775666788, 10825064499110513322, 14689361390610371514, 16734835965210899604,
	6399527266456256611, 279676139769146943, 14196707449518829319, 8999011805617471788, 6924566981437551529,
	11935101640947094708, 6946686668319032438, 1392397551035393808, 3281373847403844559, 7673207765878545335,
	11361626762965614966, 15014124176381749710, 4596876061716608039, 828591673584457147, 9455745301124654205,
	3617555776104743529, 14659471514959068984, 8574153963535421338, 14499941443938607338, 5428658603350578075,
	13021212502356346549, 5096654515527008449, 4534277910591376951, 8835565338717500304, 16576323000633271029,
	990415953277272574, 16445594914354785247, 3627100269752912500, 10426227458872807220, 1752742903139208989,
	6823688420765684666, 6032467244848876436, 10130802325065602675, 4799660975768660701, 919843791599379793,
	1400508188743108884, 9926103171667860567, 2907281439932170679, 1472519844697410736, 10494390165361325183,
	2975558351153467687, 14195456863097703305, 14488548973706263852, 13952948965625754173, 15533773800107121723,
	15246604803842169218, 14061028452424055052, 6034576862396884247, 6607332037155172840, 15095378478039910236,
	118298131398851786, 9506365343507173044, 10127547266291660615, 12627826653636866797, 7622693872122742700,
	12430169785604149026, 3175745506366470758, 11556883535617490720, 1996593920843342897, 14342363215063417557,
	12931821027969541464, 4671610853862129650, 10779167372955628948, 3056332746016649150, 16859660887416251990,
	15682387623776390888, 9757647480442399021, 9215619702456294450, 13019161915019571864, 10253388758817192048,
	16424557428031725843, 849635121368231514, 493400823683765929, 6591905403151965609, 2312873759091576466,
	17050629189067344208, 18205846881357943473, 1960528713598030433, 198614094973075395, 17941254959206722521,
	857498332500047840, 15180133523906779713, 2876636394410322752, 13359725710749045025, 13717360897469088943,
	760740741943613320, 15419901828844486474, 12914457515001554559, 17526944655607742924, 8761126201432260190,
	12168683120995651260, 4606018198686923411, 12916708273433536175, 9406074772821824226, 15172805476383067951,
	13177324915284969250, 3132227180552437724, 17408630500375936775, 2179736218039354276, 7055955377579800709,
	329778415349582010, 15934087754879685594, 14995696081579861508, 13955768992965067384, 9891590185009426703,
	15903206082924732236, 8115136352186866059, 18241648352573844123, 17040182257053926107, 7675671191860293954,
	11509334281505313682, 17647991412432377166, 692096105679558205, 16925845084292965048, 5804560326627778270,
	6933583034365165052, 15052574835258590597, 894060311800635659, 9360130382033933288, 5384925032239196514,
	9857536005042858403, 9207450753580197605, 15211078325158705070, 9089315661248809866, 16923096141865953495,
	12286048852543408741, 6100275367158842890, 18041105950862327045, 12087323835742844336, 10924076007769315270,
	18143242942063673131, 7685299261280486864, 12078452559500257731, 7432855125997443447, 13072523539970417585,
	6627273654189079311, 12974856637904390687, 8262326810234485412, 18060989106058760673, 9395971037433762449,
	7301888237939937549, 3906588315672313728, 140022567823035473, 13214308884796288464, 17568459881417548670,
	7892480169203193667, 395882274225087444, 15301855826431995565, 12431805408027574535, 13867218017115644845,
	7747147556852822068, 213148147061592444, 1147050999855880703, 11632291939583621197, 16924403520488491365,
	2903561752088116009, 9294218661857187759, 17522102186554386007, 12638487278707521906, 10147548466419665037,
	15155882069514784009, 7100973504029541625, 16520215335446872211, 11643417985196903301, 3199254614884960538,
	11203203112869441632, 16805840676370095801, 3238642280712712661, 5526591513206842612, 14557720847513007339,
	3814611319475495699, 10328797869272257725, 8408152371635163252, 496084415533370556, 2111392068471983631,
	2989977065976802023, 12458609306699398175, 3221620069009153834, 12025452106090456275, 18158478279493886774,
	6193739207526038143, 1917936835493074776, 18106765057498694960, 12010696972289491082, 7451941173504797137,
	10460080154150440624, 8741545748998466867, 15134306707821984892, 18015632174694489360, 16693143325703395845,
	10898251974987270416, 8695401978176517617, 16890217359975369936, 13785993735426160689, 7888845532527976373,
	15778827558491823086, 2522543887406335640, 3759749631911308224, 17614438810254750265, 8710526160774049443,
	13822444649968069456, 2662218518393240646, 18224417056552757317, 4174355002817818100, 414965952104828438,
	2908700881822148014, 8779784548283578595, 6931945918825248378, 4941799717679293407, 6177026811335423860,
	17203202915628020982, 13455834077139283691, 5336190215977042753, 4441328693759828549, 14301906058981948558,
	8767951844524310712, 4702678461749459233, 16071861631370481113, 13750147825259840487, 610139110116695364,
	2421049037820562398, 17759307193269279197, 12649731631824034304, 12999599482447234547, 3834635091911834437,
	10120451956124535495, 8316047019281803308, 9264366787750207329, 2526507973352074175, 10231548038441433308,
	17113491007303448887, 11594289225202895019, 18346505315258593286, 17768009244547818084, 1576698852175362199,
	12026236046955395040, 4920221184931143134, 1627829365734183404, 16727422438371618620, 15727114066671511698,
	2372320992788255040, 2309245699914256770, 12314072050221138134, 16152343871536323744, 45008050450584446,
	7756793223429037911, 2117442618385149471, 9345909586259657321, 8761626118042981173, 8295237617862628623,
	18102013175684963755, 9619980372884483533, 15796524445606297825, 471259877123105082, 12982192256927946711,
	4387533857321374987, 12428202232618429866, 9598947180905142823, 14734947463553133712, 2399226640606781105,
	16905111445918008179, 8859025831836355027, 8639602397762776516, 2115105874780061517, 15867087742919550839,
	6977317916035425249, 4151937124825921017, 12349886602107302616, 4357969982369059764, 18317564845450656092,
	12565290051078153279, 11526033050079744246, 441507584051665729, 7283855682742174347, 16198179877985385017,
	15287453646770213028, 11261726668897694944, 15287788739456521123, 12828257316277850729, 8478626249766814413,
	10603976922993760297, 4799769224472172929, 16284679893111365059, 2289079163147780107, 4200729165879719289,
	15677787644814386571, 756034889477806821, 4724875543908344324, 11680310421011052935, 5736952483423015425,
	10137133122696116249, 10519950214379636758, 10691765743379995749, 8707002453519502849, 5824343430521742329,
	11490402676837613754, 18147032868520679408, 4952771997674410968, 3119170057034093787, 15418412571793227827,
	16947631822772541382, 6768616184571698394, 4632169361618537662, 13129228254492538687, 4270785598083309515,
	18013657436889614586, 16086230615715410434, 13421754990226577623, 2291242245054327208, 15863144406366093551,
	13131335055250602130, 16072563785557161228, 6058355861822406858, 15518522449948882238, 5873359330457121544,
	11516144010999842440, 12791538999180417030, 16325100592142956067, 1165505497730535314, 2643766751137561974,
	13939643514151152253, 1831406912966687028, 9173624551887931713, 5865633732013421583, 15931705701629082341,
	3802938727468839144, 13964402647532625332, 1210013397032407120, 14255668113152614846, 7794072341422149130,
	11928442744569509227, 16401957702808161063, 10531050368304459894, 15258957834544037070, 5823159652959410424,
	12857060059065565484, 9236391232542691300, 5219198308238193528, 17411808143779738503, 4173417523730213087,
	14722149604846595111, 5929109472848581954, 16877280361542107156, 7474437698427904103, 9831513521335428760,
	751471516496209043, 12317979580658272009, 13748678753612609692, 5709541018789381327, 5789130033983192302,
	12217685426617205775, 16539484675871338433, 776486965601631849, 2345740734572741243, 6341697222245901215,
	13099193136827037181, 204123665542931154, 11224814184770612725, 721689523616564767, 14423918101495076826,
	2108234040388692286, 3390267755814225122, 11732543150467401498, 8303799689701838272, 6104304814763550045,
	4105598755364741699, 18158463576060882084, 15381018248032238483, 12584369388988336742, 61122968712918070,
	11125003388483346854, 16440246104439096383, 7069080742086533619, 16678085212164438380, 11878868338535236266,
	14849165207942182956, 13394217103534608279, 11312554262368158333, 12811783880407929025, 11500296095873059102,
	9292809724549422678, 1415563559236738350, 41406497841832993, 17609294611265847103, 2193695284893518467,
	3018469034978866138, 16538950594574727039, 13192113205170993889, 4059157342224776055, 918735250537819040,
	7231314384566776724, 4815329713662711142, 10165888923003306983, 17047438869301794399, 10734699398209093009,
	14438305441371157479, 10341221266333524138, 2539696960201429253, 9207997407084704522, 17794650224512178850,
	10176041990135375541, 11168314613925255493, 2901581893313933971, 4357389817891565749, 1321403954333962126,
	3463177043058966924, 6343871253887740107, 16490476177717155597, 9526045048457779229, 10112404192304749338,
	3088154659918113028, 14048186459060223040, 12438455572829625834, 10886085785452671065, 17660003554620478274,
	12618505071070217880, 1461491922069311906, 839491542014512762, 11650485177525672432, 12861869000360803561,
	10966171941511968938, 16194887541737579488, 15138906790616220116, 12162698788625761887, 15019308984908505957,
	2568199828224208470, 3108138781323191767, 14549575412041042417, 13173328452196668130, 1019092597011251431,
	3198658516710343231, 12737326215584292059, 16648122382433173763, 17367949995848836706, 3479918902331642405,
	1658785880069558655, 4125643631056849026, 10970024782521998339, 911755937249800616, 10741863525248742672,
	9468019022447353859, 15229392199321629467, 7226258822720332319, 1694881465938077493, 18217405464380897173,
	6744481196927712577, 4082681671017498068, 11540832052772557423, 14872104057592984841, 9430559953195483862,
	889988161859419353, 13835722178018016660, 11420300564283952476, 2591857096295729149, 4039909671596872526,
	17667212201888166861, 1895058517391226750, 6483244968225008309, 4336635309391699200, 17908815611234327373,
	10049396823908162165, 7095128839799372992, 169747289235870461, 10089654841270408837, 8769843336475300403,
	583103985801850551, 4641792461632577654, 16581270335019315057, 14253675639103960382, 3733038549661010961,
	1866003358406141496, 2934588988406802098, 17996389278185227248, 16520968670258952483, 18120880889239914402,
	10754715610447700577, 11728227432073157542, 5991972744052607848, 11583366104904301565, 5685897123094948608,
	1745774676205482621, 8858389333966407344, 12511851535884591180, 8585225526406214932, 5676097357616284031,
	6230576165160159969, 16440122771094883806, 607159736032795026, 16272455908787979999, 15902285909488726615,
	1136060433075388333, 7460344659616510194, 4934884593954878237, 14254950350147983174, 11239322631325087480,
	466438168653593478, 5792023559445812070, 5708703248236612710, 13448903291262685451, 4356206919120802411,
	17219446193933722966, 14453381737028716361, 12971651047783664417, 5608008025391549343, 13499677972396726669,
	17104440039833779066, 14970737278581341247, 6408088415333650543, 12287285914260179483, 15603913831383705565,
	7228753759322171863, 16499066678092484615, 7236572357870097629, 7858461979939016923, 12077909498898071485,
	250128204320245450, 486655662300159908, 10866479063490743369, 12747558495262780627, 16359763950954975299,
	2691316960514504584, 3441144264499340017, 3965655031128134722, 1993767464636184858, 12221718333972142595,
	4623376893212409319, 3808326428566066479, 11880397035612984008, 16351282810767720516, 2230489124048464167,
	10497807383109907535, 5298671117893846545, 12307812199765546492, 9731997401972388704, 5917492456411459044,
	10867185911947089262, 14211765313341181059, 6024333332686770941, 13819666869193803150, 11540365805209986721,
	638339116509838610, 6045154541634034788, 13045089538721505001, 11615208552883179740, 17822400142282473725,
	247165191479176190, 15789762443725506097, 17036005863388504674, 10612513721745333143, 1364926757809785840,
	5648861737609083209, 6814194137664409576, 5627117773394945862, 17160938468922591940, 9999718183463970230,
	13051790785321408270, 7378213487126013125, 1948598532820442175, 10178843542557750804, 11184450785721445417,
	4595423020975487537, 10724231774109340038, 8707652770644473705, 16788547254614540592, 16464780421016644813,
	11480553201219687790, 5659372494037539494, 2344626342596985152, 14900217313090471971, 11107963287568347231,
	9889876019586549186, 14069582549561416227, 17621329231310545341, 4294070857878064670, 7937705608936377574,
	16450355363477801970, 16913788503210988758, 17505517185107888708, 14555769404438895601, 11937068301633337067,
	726787128358804812, 9016747369828896466, 14300095693627422508, 3912702130059322190, 13638226304942148503,
	524182878498794900, 1663801210886052001, 2220702033071312548, 6070833744174116745, 3841836228334081793,
	7119288882711911681, 10415099961057165245, 11090727455094290163, 10156531605651775809, 16214263622210763423,
	7471037767326702542, 4740490797942876174, 10936719202175357162, 9338489296856347388, 4983283866355038276,
	18234815911645705457, 16558681402556737668, 13713487622875696729, 14210919335345176015, 8780278128209122769,
	6773667685205279792, 13547049726045385058, 12222686919415667035, 13136736198834712427, 4986354608351969003,
	14218210113855665316, 5889631051507738416, 16108519164677307795, 3222092199456075933, 7321507023883975762,
	1711910135236400099, 15147645464306449438, 2746396210591492110, 15765595692877958753, 16914855642374880048,
	14577270427197066839, 15165286319914035532, 9406685810087935831, 8430412867866723143, 6666894565697208155,
	15276428409540239187, 7774399337923319318, 916165650892696148, 18195847652354045630, 4738401576134308105,
	1613635449778561413, 17298012524249835796, 15414334054772349400, 107084881033925917, 7017368025567137622,
	3712026535630657102, 10056542592841047261, 12927993235529234711, 10138066000213618135, 9842521827416933143,
	1733935150091347568, 13504213180587716535, 1461631157456026954, 2604362486441655805, 900783470661715446,
	8541600586783944355, 8725731303816211947, 14854426689688703953, 10597877997330334645, 9436022003531320515,
	6818759154284360890, 9898813606805499132, 8300156826005676704, 4907094103263926114, 7450140421633622597,
	6648738534997005833, 14257929015596205723, 2311993416292370253, 12524834378486039182, 8219753787156836038,
	7373195785098309070, 11681851214347518755, 13348905036141968066, 10407059890933247482, 16372213249280201862,
	12162482843846649880, 12943058010655660340, 5219776352469082857, 663172221523735513, 2673644565579601470,
	5596029706218078403, 7620914220791404896, 15030293774579018122, 1486945396868222946, 1937190242671998327,
	15920743808333505835, 14610778576708757797, 17745961382365687969, 4312812164427198438, 6837272077571506036,
	4635637507159569590, 11838918807549355695, 8716289307662541930, 2847257467566504885, 16008715172658253473,
	12637069667802701229, 10380596670949841970, 70757813410974498, 9351894521267445977, 4338128316479634658,
	14199217252763557787, 6921102001285209759, 13622609742304416881, 7939329696646903796, 10662600983029582431,
	14219107372514411424, 2477080916348470461, 8741248032194605601, 7174218026911131881, 16024847054756727942,
	5508906111153315027, 9899389854870184465, 14196417512296220681, 6330885697262780955, 11699212642250677218,
	381114526762976423, 5977161299719085799, 3809697317681224415, 17592873634949244249, 674595638927158918,
	6479049701777514107, 10803346122259227319, 8198325534463923292, 17583155374634219918, 15024296478185921974,
	5819697006064706810, 9370552657798371447, 11279995672176659239, 5945728535650538127, 4774777776697056273,
	10019933063355970855, 2229920310511211495, 16364142987560648914, 15942562997118009950, 14729782489818652319,
	2333686872608334148, 7976782128534925663, 13897445676639730929, 8427918881758812850, 17106136842549931820,
	15572560629473818175, 2812987377548101079, 11748120390068184498, 16149672191067005055, 6689490716952148421,
	12257071124978714797, 2566986775330172491, 8466342204829421918, 8268403247082337415, 16720245021132835385,
	5567381387307206888, 8879455105364826351, 10360312203850090272, 14840859996310532183, 11196915798641287510,
	17610634310534436305, 7470195633983163316, 6976025053470300335, 14661303926120014305, 11868249999896193015,
	759537555139344696, 10705081912612298619, 13578377292298142275, 2141656950430570065, 5713501086688851293,
	6496912870202614099, 4005699469812308161, 13423244914469106099, 6443724312221752413, 17587173770950764509,
	12239451675616416223, 14750115355566919748, 359103587867291008, 10430279094932576092, 6535702772461344585,
	17573371872662400143, 5955754742858096595, 7948460627130319873, 7567411672683135103, 3865494821183132665,
	17352121233433516011, 16648157906147645989, 15559410031990244043, 15578542796594267151, 17692024017741429022,
	51498487280955413, 1718984799195005066, 17314289291872222368, 8986493835396485175, 4932145576506952883,
	6904172830867021099, 14398677982732163114, 147146823650955972, 1240618992312018422, 1273942673392086833,
	815280390251002288, 7570196185723352081, 18271648706424109628, 2168065039813320226, 146676646089303252,
	4430431962070682835, 1752755595058026148, 8373100031395864307, 11432862709512971451, 128706898611773293,
	3534543132113072466, 240166716630029087, 4673226235249861395, 13318778871388926936, 2251799714831884799,
	7219646697329593081, 5793384269586145273, 18009316682540287253, 13020351270516243175, 9058420749805154769,
	11475506936310745141, 18442523562746438147, 1293193032669976884, 13529404430340329908, 13989438296184386760,
	10198688335504589041, 14234632821298631588, 14449066687032134735, 9418123664793417205, 12977100460459283166,
	17986889256400798408, 17792235364218603600, 16578858984200943091, 13702999268618459881, 2679208585992997234,
	3030590483438781607, 6467675707017730085, 10132878602334096552, 10351540199358542442, 4439311695325339073,
	930521524945234733, 17771077579039287304, 12420767221284532628, 12364353904881135851, 1581201045240495813,
	14948198015267676722, 14013478634469585241, 2858158668136429296, 12911470040257290845, 1541854488149995751,
	11234160743996547880, 8573351219346422284, 15131527354610542084, 4628098352713031193, 5559133735586581164,
	1665891513435251620, 16476916165265834170, 4687455260814190655, 14810097204517692320, 15697411915295599631,
	15902007595972612345, 1574440850840024752, 10310026601724632614, 312856989682507124, 10781116479722575276,
	15211116982934660836, 14600777454889480047, 1247117451352823362, 2844229746743954283, 18362656956088553787,
	4614977737960852065, 12207939747921950813, 13293495406917560642, 10144163251383278433, 17460325656492349826,
	2691961910194417630, 606281927392511840, 8229370085655456780, 1100747041620789931, 1567303068252756452,
	7378930844937872259, 12555758246258556914, 14672630994419219986, 16330063599163057569, 13380093195705652690,
}
