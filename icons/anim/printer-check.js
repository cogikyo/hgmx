import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const sheet = el.querySelector('path[d^="M13.5 22"]');
  animate(sheet, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m16 19"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
