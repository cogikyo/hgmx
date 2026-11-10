import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M12.83"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const under = [el.querySelector("g"), el.querySelector('path[d^="M2.003 16.995"]')];
  animate(under, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const minus = el.querySelector('path[d="M16 17h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
