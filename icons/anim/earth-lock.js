import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('path[d^="M12 2a"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const land = el.querySelectorAll('path[d^="M7 3.34"], path[d^="M11 21.95"], path[d^="M21.54"]');
  animate(land, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const shackle = el.querySelector('path[d^="M20 6V4"]');
  animate(shackle, { translateY: [0, 1.5, 0] }, { duration: τmid, delay: τmin });
}
