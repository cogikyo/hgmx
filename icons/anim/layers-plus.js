import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M12.83"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const under = el.querySelectorAll('path[d^="M2 12"], path[d^="M2 17"]');
  animate(under, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const plus = el.querySelectorAll('path[d="M16 17h6"], path[d="M19 14v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
