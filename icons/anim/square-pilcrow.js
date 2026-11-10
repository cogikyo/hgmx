import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bowl = el.querySelector('path[d^="M12 12"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const stems = el.querySelectorAll('path[d$="v10"]');
  animate(stems, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
