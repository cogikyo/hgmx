import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const globe = el.querySelectorAll('path[d^="M12 22a"], path[d^="M15.86 14"], path[d^="M2 12h"]');
  animate(globe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m15.172 "]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
