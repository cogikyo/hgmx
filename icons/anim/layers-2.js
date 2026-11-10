import { animate, τmid } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M13 13.74"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const under = el.querySelector('path[d^="m20 14.285"]');
  animate(under, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3 });
}
