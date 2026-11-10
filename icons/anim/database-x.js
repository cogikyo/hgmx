import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M3 "], path[d^="M21"], path[d^="M19.323"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });

  const bars = el.querySelectorAll('path[d^="m17 17"], path[d^="m22 17"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: stagger(τmin, { startDelay: τmin }) });
}
