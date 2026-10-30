import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const binaural = el.querySelector('path[d^="M5 3H4"]');
  animate(binaural, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const tube = el.querySelector('path[d^="M8 15"]');
  animate(tube, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const chest = el.querySelector("circle");
  animate(chest, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
