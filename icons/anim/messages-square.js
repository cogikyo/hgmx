import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const first = el.querySelector('path[d^="M16 10"]');
  animate(first, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const reply = el.querySelector('path[d^="M20 9"]');
  animate(reply, { translateX: [-4, 1, 0], translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
