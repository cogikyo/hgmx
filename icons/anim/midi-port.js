import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const pins = ["M6 12", "M8 16", "M12 18", "M16 16", "M18 12"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(pins, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
