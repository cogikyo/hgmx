import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M8 10.37"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const under = el.querySelector("g");
  animate(under, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3 });

  const arrow = el.querySelectorAll('path[d^="M12 7"], path[d^="m9 19"]');
  animate(arrow, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
