import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const play = el.querySelectorAll('path[d^="m10.215"], path[d^="m16.042"]');
  animate(play, { translateX: [0, 2, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
