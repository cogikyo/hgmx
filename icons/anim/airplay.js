import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M5 17"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cast = el.querySelector('path[d^="m12 15"]');
  animate(cast, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
