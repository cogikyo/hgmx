import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const owl = el.querySelector("g");
  animate(owl, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const eyes = el.querySelectorAll('path[d$="9h.01"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
