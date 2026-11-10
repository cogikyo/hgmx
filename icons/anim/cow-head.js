import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const muzzle = el.querySelector("g");
  animate(muzzle, { translateX: [0, 1.2, -1.2, 0.6, 0] }, { duration: τmax });
  const eyes = el.querySelectorAll('path[d$="v.5"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
