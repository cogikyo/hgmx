import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pipette = el.querySelectorAll('path:not([d^="m2 22"])');
  animate(pipette, { translateX: [0, -2, 0.5, 0], translateY: [0, 2, -0.5, 0] }, { duration: τmid });

  const drop = el.querySelector('path[d^="m2 22"]');
  animate(drop, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
