import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const hands = el.querySelector('path[d="M12 10v2l1 1"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
  const face = el.querySelector("circle");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });
  const ring = [el.querySelector('path[d^="M2.5 9"]'), el.querySelector('path[d^="M21.5 15"]')];
  animate(ring, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
