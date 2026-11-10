import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelectorAll('path[d^="M12 4"], path[d^="M2 9"], path[d^="M16 18.2"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const soup = el.querySelector('path[d^="M16 15.6"]');
  animate(soup, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
