import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lid = el.querySelector("g");
  animate(lid, { translateY: [0, -3, 0], rotate: [0, -15, 0] }, { duration: τmax });

  const slats = el.querySelectorAll('path[d$="v6"]');
  animate(slats, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
