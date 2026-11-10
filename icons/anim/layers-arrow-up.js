import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M7.674"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const under = el.querySelector('path[d^="M2 17.002"]');
  animate(under, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3 });

  const arrow = el.querySelectorAll('path[d^="M12 12"], path[d^="m9 5"]');
  animate(arrow, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
