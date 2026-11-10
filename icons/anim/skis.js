import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const skis = el.querySelector("g.skis");
  animate(skis, { translateY: [0, -3, 1, 0] }, { duration: τmax });
  const poles = el.querySelector("g.poles");
  animate(poles, { rotate: [0, -8, 4, 0] }, { duration: τmid, delay: τmin });
}
