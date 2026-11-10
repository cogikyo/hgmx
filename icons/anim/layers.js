import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M12.83"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const under = el.querySelectorAll('path[d^="M2 12"], path[d^="M2 17"]');
  animate(under, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
