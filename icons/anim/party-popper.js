import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const popper = el.querySelector("g");
  animate(popper, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmid });

  const streamers = ["m22 2", "m22 13", "m11 2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(streamers, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const confetti = el.querySelectorAll('path[d$="h.01"]');
  animate(confetti, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
