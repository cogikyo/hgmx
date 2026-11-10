import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="H3"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const play = el.querySelector('path[d^="M15 12.003"]');
  animate(play, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
